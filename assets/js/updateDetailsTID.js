
var updateId = 0;

$(document).ready(function () {
    bindShowTIDDetails();
    bindShowUserOverride();
    bindShowFraudOverride();
    bindChangeSiteTidTableSize();
});

function bindShowTIDDetails() {
    $("[data-button='show-tid-details']").click(function () {
        var tid = $(this).attr("data-tid")

        var data = {
            TID: tid,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        }

        $.ajax({
            url: "getTidDetails",
            data: data,
            method: "POST",
            success: function (d) {
                $("#tidModalBody").html(sanitizeHTML(d))
                $("#tidModal").modal('show')
            },
            error: function (d) {
                displayWarningMessage(d.responseText)
            }
        });;
    })
}

function bindShowFraudOverride() {
    $("[data-button='override-fraud-details']").click(function () {
        var tid = $(this).attr("data-tid");

        var data = {
            TID: tid,
            ProfileId: $(this).attr("data-profile_id"),
            SiteProfileId: $(this).attr("data-profile_id"),
            SiteId: $("#siteID").val(),
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        };

        $.ajax({
            url: "/showTidFraudModal",
            data: data,
            method: "POST",
            success: function (d) {
                $("#tidFraudModalBody").html(sanitizeHTML(d));
                $("#tidFraudModal").modal('show');
                getTidVelocityLimits();
                numbersOnly()
            },
            error: function (d) {
            }
        });
    })
}

function bindShowUserOverride() {
    $("[data-button='override-user-details']").click(function () {
        var tid = $(this).attr("data-tid");

        var data = {
            TID: tid,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        };

        $.ajax({
            url: "/showTidUserModal",
            data: data,
            method: "POST",
            success: function (d) {
                $("#tidUsersModalBody").html(sanitizeHTML(d));
                $("#tidUsersModal").modal('show')
            },
            error: function (d) {
                displayWarningMessage(d.responseText)
            }
        });
    })
}

function searchUpdatesTID(tidId, siteId) {
    var data = {
        TID: tidId,
        Site: siteId,
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "updatesTID",
        data: data,
        method: "POST",
        success: function (d) {
            finalData.splice(0, finalData.length);
            localStorage.clear()
            $("#tidUpdatesModalBody").html(sanitizeHTML(d));
            bindDatePickers();
            $("#tidUpdatesModal").modal('show');
            updateId = parseInt($("#update_table tr:last").find('td:first').attr('version-select'));
            if (isNaN(updateId)) { updateId = 0; }
            document.getElementById("applyUpdates").addEventListener("click", function (event) {
                event.preventDefault()
            });
            updateId = $('#update_table').find('tr').length
            $("#siteId").val(siteId);
            $("#applyUpdates").attr("disabled", "disabled")
            removeDatePickers()
        },
    });
}

function AddTIDUpdate(tidId) {
    let data = new FormData($("#tidUpdates_" + tidId)[0]);
    data.append("TID", tidId);
    data.append("UpdateID", ++updateId);
    data.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value);
    $.ajax({
        type: "POST",
        data: data,
        url: "/addTIDUpdate",
        processData: false,
        contentType: false,
        success: function (data) {
            finalData.splice(0, finalData.length);
            $("#update_table").append(data)
            bindDatePickers()
            $("#applyUpdates").prop("disabled", false);
            removeDatePickers()
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function ApplyTIDUpdates(tidId, tidPofileID) {
    let dataObj = new FormData($("#tidUpdates_" + tidId)[0]);
    let apkIDs = []
    let updateID = 0
    let apkID = $("#apkIDs").val();
    if (apkID != null) {
        apkIDs = apkID;
    }
    let tidUpdateID = $("#tidUpdateID").val();
    if (tidUpdateID != null) {
        updateID = tidUpdateID
    }
    dataObj.append("tidPofileID", tidPofileID);
    dataObj.append("TID", tidId);
    dataObj.append("tidUpdateID", updateID);
    dataObj.append("thirdParty_" + updateID, apkIDs);
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value);
    $.ajax({
        type: "POST",
        data: dataObj,
        url: "/ApplyTIDUpdates",
        processData: false,
        contentType: false,
        success: function (data) {
            finalData.splice(0, finalData.length);
            localStorage.clear()
            $("#update_table").append(sanitizeHTML(data))
            $("#applyUpdates").attr("disabled", "disabled")
            $("#update_table").find("input,select").attr("disabled", "disabled");
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function DeleteUpdate(tidId, updateID, siteId,profileId) {
    confirmDialog("Delete update?", 'Are you sure you want to delete update?', function () {
        var data = {
            TID: tidId,
            UpdateID: updateID,
            SiteId: siteId,
            ProfileId: profileId,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        }
        var row = $("#updateDetailsRow" + updateID);
        row.remove()

        $.ajax({
            type: "POST",
            data: data,
            url: "/DeleteTIDUpdate",
            success: function (data) {
            },
            error: function (data) {
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function removeDatePickers() {
    $(".xdsoft_datetimepicker").remove();
}

function closeThirdPartyModal() {
    finalData.splice(0, finalData.length);
    counter = 0;
}

function closeDataModal() {
    $(".xdsoft_datetimepicker").detach()
    $("#tidUpdatesModal").prependTo("body");

    finalData.splice(0, finalData.length);
    counter = 0;
    siteSearch(getCurrentPage(), getProfileId(), getProfileType(), getCurrentPageSize(), getSearchTerm());
}

var checkPastTime = function (inputDateTime) {
    let inputDateTimePicker = "#date" + updateId
    if (typeof (inputDateTime) != "undefined" && inputDateTime !== null) {
        var current = new Date();

        //check past year and month
        if (inputDateTime.getFullYear() < current.getFullYear()) {
            $(inputDateTimePicker).datetimepicker('reset');
        } else if ((inputDateTime.getFullYear() == current.getFullYear()) && (inputDateTime.getMonth() < current.getMonth())) {
            $(inputDateTimePicker).datetimepicker('reset');
        }

        // 'this' is jquery object datetimepicker
        // check input date equal to current date
        if (inputDateTime.getDate() == current.getDate()) {
            if (inputDateTime.getHours() < current.getHours()) {
                $(inputDateTimePicker).datetimepicker('reset');
            }
            this.setOptions({
                minTime: current.getHours() + current.getMinutes() + ':00' //here pass current time hour
            });
        } else {
            this.setOptions({
                minTime: false
            });
        }
    }
};
var currentYear = new Date();
function bindDatePickers() {
    $('.inputDateTime').datetimepicker({
        format: 'Y-m-d H:i',
        minDate: 0,
        step: 5,
        yearStart: currentYear.getFullYear(), // Start value for current Year selector
        onChangeDateTime: checkPastTime,
        onShow: checkPastTime
    });
}

function changeSiteTidTablePrevious(currentPage) {
    changePage(currentPage - 1, getProfileId(), getProfileType(), getCurrentPageSize(), getSearchTerm());
}

function changeSiteHistory() {
    hideWarning()
    $("#profLoader").show()
    profileChangeHistory(getProfileId(), getProfileType(), getSiteId(), 1);
}

function appendTidshtml(divId) {
    $('#'+divId).html('<div class="modal fade bd-example-modal-lg" id="tidSelectModal" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">\n' +
    '    <div class="modal-dialog modal-dialog-scrollable modal-lg" role="document">\n' +
    '        <div class="modal-content">\n' +
    '            <div class="modal-header">\n' +
    '                <h5 class="modal-title" id="update-modal-title">Select TIDs</h5>\n' +
    '            </div>\n' +
    '            <div class="modal-body">\n' +
    '                <div id="flag-status" class="add-site-data-groups" >\n' +
    '                    </div>\n' +
    '                    <button id="flag-status-update" type="button" class="btn btn-primary float-end" onclick="hideTIDSelectModal()">Close</button>\n' +
    '            </div>\n' +
    '        </div>\n' +
    '    </div>\n' +
    '</div>');
}

function clearSiteHistory(event) {
    hideWarning()
    $('#history').html('');
    $('#users_tid_data_for_flagging div').remove();
    $('#fraud_tid_data_for_flagging div').remove();
    $('#site_tid_data_for_flagging div').remove();

    if (event.id=="site-tab"){
        hideWarning()
        appendTidshtml("site_tid_data_for_flagging");
    }else if(event.id=="users-tab") {
        hideWarning()
        appendTidshtml("users_tid_data_for_flagging");
    }else if(event.id=="fraud-tab") {
        hideWarning()
        appendTidshtml("fraud_tid_data_for_flagging");
    }
}

function changeSiteTidTableNext(currentPage) {
    changePage(currentPage + 1, getProfileId(), getProfileType(), getCurrentPageSize(), getSearchTerm());
}

function changeSiteTidTablePage(targetPage) {
    if (targetPage === "First") {
        targetPage = 1;
    } else if (targetPage === "Last") {
        targetPage = $("#pageCount")[0].value;
    }
    changePage(targetPage, getProfileId(), getProfileType(), getCurrentPageSize(), getSearchTerm());
}

function changeSiteTidTableSize(targetSize) {
    changePage(1, getProfileId(), getProfileType(), targetSize, getSearchTerm());
}

function bindChangeSiteTidTableSize() {
    $('select[name="site-tid-table-length"]').change(function () {
        changeSiteTidTableSize($(this).val());
    })
}

function searchSiteTids() {
    changePage(1, getProfileId(), getProfileType(), getCurrentPageSize(), getSearchTerm());
}

function clearSiteTidSearch() {
    $('#site-tid-table-search-term').val('');
    changePage(1, getProfileId(), getProfileType(), getCurrentPageSize(), "")
}

function getProfileId() {
    return $("#profileID")[0].value;
}

function getSiteId() {
    return $("#siteID")[0].value;
}

function getProfileType() {
    return $("#profileType")[0].value;
}

function getCurrentPageSize() {
    return $("#pageSize")[0].value;
}

function getSearchTerm() {
    return $('#site-tid-table-search-term').val();
}

function getCurrentPage() {
    return $("#pageNumber")[0].value;
}

function changePage(targetPage, profileId, profileType, pageSize, searchTerm) {
    removeDatePickers()

    $("#profLoader").show()
    siteSearch(targetPage, profileId, profileType, pageSize, searchTerm);
}

function siteSearch(targetPage, profileId, profileType, pageSize, searchTerm) {
    $.ajax({
        url: '/profileMaintenance',
        data: {
            profileID: profileId,
            pageSize: pageSize,
            pageNumber: targetPage,
            pageChange: "true",
            searchTerm: searchTerm,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        method: 'POST',
        success: function (data) {
            let clean = sanitizeHTML(data)
            $('#tids').html(clean);
            $("#profLoader").hide();
            hideWarning();
            bindSubmit();
            bindOverrideField();
            bindShowFraudOverride();
            bindShowUserOverride();
            bindShowTIDDetails();
            bindChangeSiteTidTableSize();
            bindToggle();
            bindDatePickers();
            $('[data-button="toggle"]').parents().next('.hide').toggle();
            $('[name*="multi."]').multiselect()
        },
        error: function (data) {
            if (data.status === 401 || data.status === 403) {
                displayWarningMessage('User not authorised');
            } else {
                displayWarningMessage(data.responseText);
            }
        }
    });
}

function profileChangeHistory(profileId, profileType, siteId, pageNumber) {
    $.ajax({
        url: '/profileMaintenanceChangeHistory',
        data: {
            profileID: profileId,
            siteID: siteId,
            ProfileType: profileType,
            pageNumber: pageNumber,
            pageSize: 50,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        method: 'POST',
        success: function (data) {
            $("#profLoader").hide();
            let clean = sanitizeHTML(data);
            $('#history').show();
            $('#history').html(clean);
        },
        error: function (data) {
            $("#profLoader").hide();
            if (data.status === 401 || data.status === 403) {
                displayWarningMessage('User not authorised');
            } else {
                displayWarningMessage(data.responseText);
            }
        }
    });
}


function manageThirdPartyTID(tidId, siteId, tidUpdateId) {
    let btnID = event.target.id;
    let packageIDs = localStorage.getItem(btnID.toString())
    let data = {
        TID: tidId,
        Site: siteId,
        tidUpdateId: tidUpdateId,
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "updatesThirdParty",
        data: data,
        method: "POST",
        success: function (d) {
            $("#third-party-edit-modal-body").html(sanitizeHTML(d));
            $('#third-party-select option').remove();
            $("#third-party-select").append($("<option selected disabled hidden />").val("").text("Select 3rd Party Application to Add"));
            $('#addThirdParty').attr('disabled', 'disabled');
            $("#third-party-select").attr("disabled", "disabled");
            $("#third-party-edit-modal").modal('show');
            $('#btnID').val(btnID);
            addThirdPartyTarget(tidId, siteId, tidUpdateId, packageIDs)

        },
    });
}

function addThirdPartyTarget(tidId, siteId, tidUpdateId, packageIDs) {

    hideWarning("updatesThirdPartyTID");
    let thirdPartyTargetDataObj = {
        field: []
    };
    if (tidUpdateId) {
        let data = {
            TID: tidId,
            Site: siteId,
            tidUpdateId: tidUpdateId,
            packageIDs: packageIDs,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        }
        $.ajax({
            type: 'POST',
            data: data,
            url: '/getThirdPartyTarget',
            success: function (data) {
                $('#tidUpdateID').val(tidUpdateId);
                hideWarning("updatesThirdPartyTID");
                let apkIds = [];
                for (let i = 0; i < data.length; i++) {
                    for (const [key, value] of Object.entries(data[i])) {
                        if (key === "ApkID") {
                            apkIds.push(value)
                        } else {
                            finalData.push(value)
                        }
                        thirdPartyTargetDataObj.field.push({
                            name: key,
                            value: value
                        });
                    }
                }
                $("#apkIDs").val("[" + apkIds + "]");
                displayData(thirdPartyTargetDataObj.field)
            },
            error: function (data) {
                displayWarningMessage(data.responseText)
            }
        })
    } else {
        let apk = $('#third-party-select').find(':selected');
        let apkReadable = $('#third-party-select option:selected').text();
        let partialPackageName = $('#partialPackageName_select option:selected').text();

        if (apk.val() === "" || apk.val() === undefined) {
            displayWarningMessage("Please select the 3rd Party Target Package", "updatesThirdPartyTID")
            return false
        }

        if (!validateThirdPartyTarget(finalData, partialPackageName)) {
            displayWarningMessage("3rd Party Target Package already added", "updatesThirdPartyTID")
            return false
        } else {
            thirdPartyTargetDataObj.field.push({
                name: "ApkID",
                value: apk.val()
            });
            thirdPartyTargetDataObj.field.push({
                name: "ThirdPartyApkID",
                value: apkReadable
            });
            finalData.push(apkReadable)
            displayData(thirdPartyTargetDataObj.field)
        }
    }
}

function getThirdPartyList(partialPackageName, tidId, siteId, packageIDs) {
    if (partialPackageName.value === "") {
        $('#third-party-select option').remove();
        $("#third-party-select").append($("<option selected disabled hidden />").val("").text("Select 3rd Party Application to Add"));
        $("#third-party-select").attr("disabled", "disabled");
        $('#addThirdParty').attr('disabled', 'disabled');
        return false
    }
    let data = {
        TID: tidId,
        Site: siteId,
        partialPackageName: partialPackageName.value,
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "updatesThirdParty/Select",
        data: data,
        method: "POST",
        success: function (d) {
            if (d === null) {
                $("#third-party-select").append($("<option selected disabled hidden />").val("").text("No valid third party application"));
                $("#third-party-select").attr("disabled", "disabled");
                $('#addThirdParty').attr('disabled', 'disabled');
                return
            }
            hideWarning("updatesThirdPartyTID");
            $('#third-party-select option').remove();
            $("#third-party-select").append($("<option selected disabled hidden />").val("").text("Select 3rd Party Application to Add"));
            $("#third-party-select").attr("disabled", "disabled");
            $('#addThirdParty').attr('disabled', 'disabled');
            for (let i = 0; i < d.length; i++) {
                let tpApk = d[i]
                $("#third-party-select").append($("<option />").val(tpApk["ApkID"]).text(tpApk["Apk"]));
                $("#third-party-select").removeAttr("disabled")
                $('#addThirdParty').removeAttr('disabled');
            }
        },
    });
}

let counter = 0;
let finalData = [];

let thirdParty = document.getElementById('third-party-edit-modal');
let tidUpdates = document.getElementById('tidUpdatesModal');
let snUpdates = document.getElementById('snUpdatesModal');
window.onclick = function (event) {
    if (event.target == thirdParty || event.target == tidUpdates || event.target == snUpdates) {
        finalData.splice(0, finalData.length);
        counter = 0;
    }
}

$("#third-party-edit-modal").on('hidden.bs.modal', function () {
    finalData.splice(0, finalData.length);
    counter = 0;
});

$("#tidUpdatesModal").on('hidden.bs.modal', function () {
    finalData.splice(0, finalData.length);
    counter = 0;
});

$("#snUpdatesModal").on('hidden.bs.modal', function () {
    finalData.splice(0, finalData.length);
    counter = 0;
});

function validateThirdPartyTarget(finalData, currentValue) {
    if (finalData.some(data => data.toLowerCase().includes(currentValue.toLowerCase()))) {
        return false
    }
    return true
}

function removeThirdPartyTarget(value) {
    let tpApkName = $(value).closest("tr").find(".ThirdPartyApkID").text();
    finalData = $.grep(finalData, function (value) {
        return value != tpApkName;
    });
    $(value).parents("tr").remove();
}

function displayData(data) {
    let rows = [];
    rows = data

    let list = '<tr>';
    $.each(rows, function (i, field) {
        if (field.name === "ApkID") {
            list += '<td hidden="hidden"><input type="hidden" id="' + field.name + '" name="' + field.name + '" value="' + field.value + '"/></td>';
        } else {
            list += '<td class="ThirdPartyApkID">' + field.value + '<input type="hidden" name="' + field.name + String(counter) + '" id="' + field.name + String(counter) + '" value="' + field.value + '"/>' + '</td>';
            list += '<td><button class="btn btn-warning" id="remove-third-party" onclick="removeThirdPartyTarget(this);">delete</button></tr>';
        }

    });
    $('#userList tbody').append(list);
}

function saveThirdPartyTargetData(tidId, siteId) {
    let apkIds = []
    let dataObj = new FormData($("#thirdPartyUpdate_")[0]);
    for (const value of dataObj.getAll("ApkID")) {
        apkIds.push(Number(value))
    }

    let tidUpdateID = $("#tidUpdateID").val();
    dataObj.append("TID", tidId);
    dataObj.append("Site", siteId);
    dataObj.append("tidUpdateID", tidUpdateID);

    let btnID = dataObj.get("btnID");

    storeIdsInLocalStorage(btnID, apkIds)

    let packageIds = document.createElement("input");
    packageIds.setAttribute("type", "hidden");
    packageIds.setAttribute("name", "packageIds");
    packageIds.setAttribute("id", "packageIds");
    packageIds.setAttribute("value", "[" + apkIds + "]");
    dataObj.append("packageIds", "[" + apkIds + "]");
    $("#apkIDs").val("[" + apkIds + "]");
    $.ajax({
        type: "POST",
        data: dataObj,
        url: "/updatesThirdPartyApks",
        processData: false,
        contentType: false,
        success: function (d) {
            displayWarningMessage("Saved Successfully", "updatesThirdPartyTID")
            $('#partialPackageName_select')[0].selectedIndex = 0;
            $('#third-party-select option').remove();
            $("#third-party-select").append($("<option selected disabled hidden />").val("").text("Select 3rd Party Application to Add"));
            $("#third-party-select").attr("disabled", "disabled");
            $('#addThirdParty').attr('disabled', 'disabled');
        },
    });
}

function storeIdsInLocalStorage(btnID, apkIds) {
    localStorage.removeItem(btnID.toString());
    localStorage.setItem(btnID.toString(), JSON.stringify(apkIds));
}

function tidDetailsDisplay(elementId) {
    var tbody = document.getElementById(elementId);
    if (tbody.style.display === "none") {
        tbody.style.display = "table-row-group";
    } else {
        tbody.style.display = "none";
    }
}

function changeSiteHistoryTablePrevious(currentPage) {
    profileChangeHistory(getProfileId(), getProfileType(), getSiteId(), currentPage - 1)
}

function changeSiteHistoryTablePage(currentPage) {
    profileChangeHistory(getProfileId(), getProfileType(), getSiteId(), currentPage)
}

function changeSiteHistoryTableNext(currentPage) {
    profileChangeHistory(getProfileId(), getProfileType(), getSiteId(), currentPage + 1)
}