$(document).ready(function () {
    bindOverrideField();
    bindMultiSelect();
    bindUpdateDataGroups();
});

function bindUpdateDataGroups() {

    $("#data-group-form").submit( function(e) {
        e.preventDefault()

        let profileId = $("#profileID").val();
        let profileType = $("#profileType").val().toLowerCase();

        let dataObj = new FormData($("#data-group-form")[0]);

        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: '/updateDataGroups',
            processData: false,
            contentType: false,

            success: function (data) {
                onDataGroupUpdate()
                $(location).attr('href', '/profileMaintenance?profileId=' + profileId + '&type=' + profileType + '&DGUpdated=true');

            },
            error: function (data) {
                if (data.status === 401 || data.status === 403) {
                    displayWarningMessage('User not authorised')
                } else {
                    displayWarningMessage(data.responseText)
                }
            }
        })
    })


}

function onDataGroupUpdate () {
    // create an observer to wait until siteTab is active
    let observer = new MutationObserver((mut, obs) => {
        let siteTab = $("#site-tab")
        if (siteTab.hasClass("active")) {
            // if active, bind multi-select then stop observing
            bindMultiSelect()
            obs.disconnect()
        }
    })
    observer.observe(document, {
        childList: true,
        subtree: true,
    })
}

let removedOverrides = new Set()

function toggleRemoveOverride(elementId) {
    const overrideButton = $('#override-button-' + elementId)
    // select element and the button child of the next element div if exists
    const element = $('#' + elementId + ', #' + elementId + ' + div button')
    if (!removedOverrides.has(elementId)) {
        removedOverrides.add(elementId)
        overrideButton.addClass('btn-danger')
        overrideButton[0].innerText = 'Reinstate Override'
        if (elementId === "endOfDay-time"){
            $('#endOfDayAutoType').prop('disabled',true)
            $('#endOfDay-time').prop('disabled',true)
            $('#endOfDay-time-2').prop('disabled',true)
        }else{
            element.prop('disabled', true)
        }

    } else {
        removedOverrides.delete(elementId)
        overrideButton.removeClass('btn-danger')
        overrideButton[0].innerText = 'Remove Override'
        if (elementId === "endOfDay-time"){
            $('#endOfDayAutoType').removeAttr('disabled')
            $('#endOfDay-time').removeAttr('disabled')
            $('#endOfDay-time-2').removeAttr('disabled')
        }else{
            element.removeAttr('disabled')
        }

    }
}

function togglePasswordInputField (id, originalType) {
    // select the element with the correct ID, ignoring the row
    let input = $('#' + id)
    if (input.attr('type') === 'password') {
        switch (originalType) {
            case 'INTEGER':
                input.removeAttr('readonly')
                input.attr('type', 'number')
                break
            case 'STRING':
                input.removeAttr('readonly')
                input.attr('type', 'text')
                break
        }
    } else {
        input.attr('readonly', true)
        input.attr('type', 'password')
    }
}

function togglePasswordDiv (shownid, hiddenid) {
    let shown_input = document.getElementById(shownid)
    let hidden_input = document.getElementById(hiddenid)

    shown_input.hidden = !shown_input.hidden
    hidden_input.hidden = !hidden_input.hidden
}


function bindMultiSelect() {
    $('[name*="multi."]').multiselect()
}

function bindOverrideField() {
    hideWarning()
    $("[data-button='override-field']").click(function () {
        var groupName = $(this).attr("data-group-name")
        let groupDisplayName = $(this).attr("data-group-display-name")
        var elementId = $(this).attr("data-element")
        var ElementVal = ""
        var data = {
            ProfileId: $("#profileID")[0].value,
            ElementId: elementId,
            ElementValue: ElementVal,
            Group: groupName,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        }
        $.ajax({
            url: "getElementData",
            data: data,
            method: "POST",
            success: function (d) {
                if ($("#row-" + elementId)[0] == undefined) {
                    var group = $("#group-" + groupName)
                    if (group.length !== 0) {
                        let element = $.parseHTML($.trim(d))
                        let existingObject = group[0].children.namedItem(element[0].id)
                        if (existingObject == null) {
                            let clean = sanitizeTableHTML(d)
                            group.append(clean)
                        }
                    } else {
                        var groupTable = $("#group-table")

                        $tableHeaderSection = $("<thead>", {class: "thead-light"});
                        $headerRow = $("<tr>");
                        $header = $("<th>", {id: "group-header-" + groupName, colspan: "3"});
                        $header.append(groupDisplayName);
                        $headerRow.append($header)
                        $tableHeaderSection.append($headerRow)

                        $tableBodySection = $("<tbody>", {id: "group-" + groupName})
                        let clean = sanitizeTableHTML(d)
                        $tableBodySection.append(clean)

                        groupTable.append($tableHeaderSection)
                        groupTable.append($tableBodySection)
                    }

                    $('[name="multi.' + elementId + '"]').multiselect()

                    $("#defaults-tab").removeClass("active")
                    $("#site-tab").addClass("active")
                    $("#default").removeClass("active")
                    $("#site").addClass("active")
                    $('[data-bs-toggle="tooltip"]').tooltip();
                }
            },
            error: function (data) {
                displayWarningMessage(data.responseText)
            }
        });
        ;
    })
}

function hideTIDSelectModal(){
    $("#tidSelectModal").modal('hide');
}

function showTIDSelectModal(){
    $("#tidSelectModal").modal('show');
}

function getFlagStatus(){
    showTIDSelectModal()
    var data = {
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
        siteID : $("#siteID")[0].value,
    }

    $.ajax({
        url: "getFlagStatus",
        data: data,
        method: "POST",
        success: function (d) {
            $('#flag-status').html(sanitizeHTML(d));
        },
        error: function(data){
            //the data is a json array of messages
            if(data.getResponseHeader("content-type") === "application/json") {
                let validationMessages = JSON.parse(data.responseText).join('<br/>');
                displayWarningMessage(validationMessages)
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    });
}

function flagStatusSelect(flagging, flagStatusId) {
    switch (flagging.value) {
        case "all":
            $("#"+flagStatusId).html('')
            $("#hidden_"+flagStatusId).val(flagging.value);
            break
        case "specific":
            var d = '<button type="button" class="btn btn-info p-1" onclick="showTIDSelectModal()"> View / Update </button>'
            $("#"+flagStatusId).html(d);
            $("#hidden_"+flagStatusId).val(flagging.value);
            getFlagStatus()
            break
        case "file":
            var d = '<input id="flagged-tids-file" name="flagged-tids-file" type="file" className="form-control-file"/>'
            $("#"+flagStatusId).html(d);
            $("#hidden_"+flagStatusId).val(flagging.value);
            break
    }
}


function downloadRpiCertificate() {
    let req = new XMLHttpRequest();

    var alphaRunes = "abcdefghijklmnopqrstuvwxyz";
    var alphaCapRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    var numericRunes = "0123456789";
    var specialRunes = "!\"#$%^&*";
    var allRunes = alphaRunes + alphaCapRunes + numericRunes + specialRunes;

    // Try finding the MID from the value held in the config item store/merchantNo
    var mid = getTextIfExists('store-merchantNo');

    // If we fail to find it there we may be creating a new site which, on rare occasions, sometimes does not have the 'store'
    // data group auto populated; if so, try to get the MID from the 'MID' text box (which only appears during site creation)
    if (mid === "") {
        mid = getTextIfExists('mid')
    }

    var password = "";
    password += alphaRunes.charAt(Math.floor(Math.random() * alphaRunes.length));
    password += alphaCapRunes.charAt(Math.floor(Math.random() * alphaCapRunes.length));
    password += numericRunes.charAt(Math.floor(Math.random() * numericRunes.length));
    password += specialRunes.charAt(Math.floor(Math.random() * specialRunes.length));
    while (password.length < 8) {
        password += allRunes.charAt(Math.floor(Math.random() * allRunes.length));
    }
    var pw = password.split("");
    for (var i = pw.length - 1; i > 0; i--) {
        var j = Math.floor(Math.random() * i + 1);
        var tmp = pw[i];
        pw[i] = pw[j];
        pw[j] = tmp;
    }
    password = pw.join("");
    var passwordb = btoa(password);

    let params = "MID=" + mid + "&PASSWORD=" + passwordb;
    req.open("Get", "/downloadRpiCertificate?" + params, true);
    req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    req.responseType = "blob";
    req.onload = function (event) {
        var blob = req.response;
        var fileName = req.getResponseHeader("fileName"); //if you have the fileName header available
        if (fileName != null) {
            var link = document.createElement('a');
            link.href = encodeURI(window.URL.createObjectURL(blob));
            link.download = fileName;
            link.click();
        } else {
            var reader = new FileReader();
            reader.onload = function () {
                displayWarningMessage(reader.result);
            };
            reader.readAsText(blob);
        }
        alert("Certificate password: " + password)
    };
    req.send(params);
}

function getTextIfExists(id) {
    var selector = $('#' + id);
    if (selector && selector.length && typeof(selector[0].value) !== 'undefined') {
        return selector[0].value
    } else {
        return ""
    }
}
