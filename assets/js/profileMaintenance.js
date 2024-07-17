var tidId = 0


let paymentServices = [];
let paymentServicesTemp = [];
let serviceInputElementId;
const regexTid = /^(?!0{8})[0-9]{8}$/
const regexMid = /^(?!0+$)[a-zA-Z0-9]{6,15}$/

var otpIntentEnum = {
    Enrolment: 0,
    Reset: 1
}

$(document).ready(function () {
    $('[data-bs-toggle="tooltip"]').tooltip();
});

function bindSaveSiteFraudLimits() {
    $('#saveSiteFraudLimits').submit(function (e) {
        e.preventDefault()
        const dataObj = new FormData(document.querySelector("#saveSiteFraudLimits"));
        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            processData: false,
            contentType: false,

            success: function (data) {
                window.location.href = "/search"
            },
            error: function (data) {
                //the data is a json array of messages
                if (data.getResponseHeader("content-type") === "application/json") {
                    let validationMessages = JSON.parse(data.responseText).join('<br/>');
                    displayWarningMessage(validationMessages)
                } else {
                    displayWarningMessage(data.responseText)
                }
            }
        })
    })
}

function bindSubmit() {
    let saveProfile = $('#saveProfile')
    saveProfile.off('submit') // Remove existing handlers
    saveProfile.submit(function (e) {
        e.preventDefault()

        const saveButton = $('#site-save')[0]
        // We disable the button when the form is submitted to stop users from pressing it multiple times which would
        // cause multiple site creation requests to be sent. We don't need to enable it again in the success of the ajax
        // request as we direct to  the search page.
        saveButton.disabled = true

        var dataObj = new FormData(document.querySelector('#saveProfile'))
        let currentOverrides = dataObj.get('removeOverrides')
        if (currentOverrides == null) {
            currentOverrides = ''
        }
        for (let id of removedOverrides) {
            currentOverrides += id + ','
        }
        dataObj.append('removeOverrides', currentOverrides)

        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr('action'),
            processData: false,
            contentType: false,

            success: function (data) {
                window.location.href = '/search'
            },
            error: function (data) {
                saveButton.disabled = false
                displayWarningMessage(data.responseText)
            }
        })
    })
}

$(document).ready(function () {
    $('[id^=multiselect]').multiselect(); //ID changed from multi to multiselect, as this was causing NEX-6379
    bindSubmit();
    bindSaveSiteFraudLimits();
    $('#data-group-form').submit(function (e) {
        e.preventDefault();

        $.ajax({
            data: $(this).serialize(),
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            success: function (d) {
                let clean = sanitizeHTML(d);
                $("#site").html(clean)
                bindSubmit()
                bindOverrideField()
                hideWarning()
                $('[id^=multiselect]').multiselect();
                $("#data-groups-tab").removeClass("active")
                $("#site-tab").addClass("active")
                $("#data-groups").removeClass("active")
                $("#site").addClass("active")
            },
            error: function (data) {
                displayWarningMessage(data.responseText)
            }

        })
    })
    bindToggle();
    $('[data-button="toggle"]').parents().next('.hide').css('display', 'none')
})

function bindToggle() {
    $('[data-button="toggle"]').click(function () {
        $(this).parents().next('.hide').toggle();
    });
}

function AddTidOnClick() {
    let data = {
        TidIndex: tidId,
        Site: $("#siteID")[0].value,
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }

    $.ajax({
        url: "addNewTID",
        data: data,
        method: "POST",
        success: function (d) {
            let clean = sanitizeTableHTML(d);
            $("#tid_table").prepend(clean)
            $('[name*="multi."]').multiselect();
        },
        error: function (data) {
            displayWarningMessage("Error creating TID: " + data.responseText)
        }
    });
    ++tidId;
}

function deleteTID(tidId) {
    confirmDialog('Delete TID?', 'Are you sure you want to delete TID: ' + tidId.toString().padStart(8, '0') + ' ?', () => {
        const data = {
            TID: tidId.toString().padStart(8, '0'),
            Site: $('#siteID')[0].value,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        }

        $.ajax({
            url: 'deleteTID',
            data: data,
            method: 'POST',
            success: (d) => {
                $('#tid-row-' + tidId).remove()
                hideWarning()
                $('[id^=multiselect]').multiselect()
                bindSubmit()
                bindOverrideField()
                $('#tid-override-' + tidId).remove()
            },
            error: (data) => {
            }
        })
    })
}

function searchDeleteTID(tidId, siteID, success) {
    confirmDialog('Delete TID?', 'Are you sure you want to delete TID: ' + tidId.toString().padStart(8, '0') + ' ?', function () {
        $.ajax({
            url: 'deleteTID',
            data: {
                TID: tidId.toString().padStart(8, '0'),
                Site: siteID != null ? siteID : $('#siteID' + tidId)[0].value,
                csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
            },
            method: 'POST',
            success: success || function (d) {
                $('#tid-row-' + tidId).remove()
            },
            error: (d) => {
                alert(d.responseText)
            }
        })
    })
}

function ApplyNewTID(tidId) {
    $(".xdsoft_datetimepicker").remove();
    let dataObj = new FormData($("#addTidProfile_" + tidId)[0]);
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value)

    $.ajax({
        url: "addTID",
        data: dataObj,
        method: "POST",
        processData: false,
        contentType: false,
        success: function (d) {
            let clean = sanitizeHTML(d);
            $('#tids').html(clean);
            hideWarning();
            $('[id^=multiselect]').multiselect();
            bindSubmit();
            bindOverrideField();
            bindShowFraudOverride();
            bindShowUserOverride();
            bindShowTIDDetails();
            bindToggle();
            $('[data-button="toggle"]').parents().next('.hide').toggle();
        },
        error: function (data) {
            //the data is a json array of messages
            if (data.getResponseHeader("content-type") === "application/json") {
                let validationMessages = JSON.parse(data.responseText).join('<br/>');
                displayWarningMessage("Error saving TID: " + validationMessages)
            } else {
                displayWarningMessage("Error saving TID: " + data.responseText)
            }
        }
    });
}

function GenerateOTP(tidId, intent, success) {
    $.ajax({
        url: "generateOTP",
        data: {
            tid: tidId,
            intent: intent,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
        },
        method: "POST",
        success: success || function (d) {
            switch (intent) {
                case otpIntentEnum.Enrolment:
                    var pinCell = $("#EnrolmentPinCell" + tidId)
                    pinCell[0].innerText = d.PIN
                    break;
                case otpIntentEnum.Reset:
                    var pinCell = $("#ResetPinCell" + tidId)
                    pinCell[0].innerText = d.PIN
                    break;
                default:
                    break;
            }
        }
    })
}

function clearInput(elementId) {
    $('#' + elementId).val('')
    $('#' + elementId + '-image').attr('src', '')
}

function openChooseImageFile(elementId) {

    var postData = {
        ButtonText: "Select",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "getFileList?type=image",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindChooseFile(elementId, "")
            $("#chooseFileModal").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}
function openChooseSoftUIConfigFile(elementId) {
    var postData = {
        ButtonText: "Select",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "getFileList?type=softUIConfig",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindChooseFile(elementId, "softUIConfig")
            $("#chooseFileModal").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function openChooseMenuFile(elementId) {
    var postData = {
        ButtonText: "Select",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "getFileList?type=menu",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindChooseFile(elementId, "menu")
            $("#chooseFileModal").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function openChooseReceiptFile(elementId) {
    var postData = {
        ButtonText: "Select",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "getFileList?type=receiptConfig",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindChooseFile(elementId, "receiptconfig")
            $("#chooseFileModal").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function openChooseTextFile(elementId) {

    var postData = {
        ButtonText: "Select",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "getFileList?type=text",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindChooseFile(elementId, "")
            $("#chooseFileModal").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

let modalVars;
let jsonEditModalOutputFieldId;
let configModalCheck;
function showJsonEditModal(modalTemplateEndpoint, elementName, outputFieldId) {
    configModalCheck = 1
    $.ajax({
        url: modalTemplateEndpoint,
        data: {
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            jsonEditModalOutputFieldId = outputFieldId;
            $("#config-edit-modal-title").text(elementName);
            $("#config-edit-modal-body").html(clean);
            $("#config-edit-modal").show()
            $("#config-edit-modal").modal('show');
            $(".modal-backdrop").show();
            renderCashbackTable();
            let definitionsJson = document.getElementById(outputFieldId).getAttribute("value")
            let definitionsObj = JSON.parse(JSON.stringify(JSON.parse(definitionsJson) , ["BIN","MIN Purchase Amount","MAX Cashback Amount","CheckType"]));
            let editRows = $(modalVars.cashbackTableID + ' tbody tr');
            editRows.each((i) => {
                document.getElementById("bin-row-" + i).innerHTML = definitionsObj[i]["BIN"];
                document.getElementById("min-purchase-amount-row-" + i).innerHTML = definitionsObj[i]["MIN Purchase Amount"];
                document.getElementById("max-cashback-amount-row-" + i).innerHTML = definitionsObj[i]["MAX Cashback Amount"];
                document.getElementById("check-type-row-" + i).value = definitionsObj[i]["CheckType"];
            });
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function closeJsonEditModal() {
    jsonEditModalOutputFieldId = "";
    $("#config-edit-modal").modal('hide');
}

let modalDpoVars;
let jsonEditModalDpoOutputFieldId;
let dpoModalCheck;

function showJsonEditDpoModal(modalTemplateEndpoint, elementName, outputFieldId) {
    dpoModalCheck = 1
    $.ajax({
        url: modalTemplateEndpoint,
        data: {
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            jsonEditModalDpoOutputFieldId = outputFieldId;
            $("#dpo-config-edit-modal-title").text(elementName);
            $("#dpo-config-edit-modal-body").html(clean);
            $("#dpo-config-edit-modal").show()
            $("#dpo-config-edit-modal").modal('show');
            $(".modal-backdrop").show();
            renderDpoTable();
            let definitionsDpoJson = document.getElementById(outputFieldId).getAttribute("value")
            let definitionsDpoObj = JSON.parse(definitionsDpoJson)
            let editRows = $(modalDpoVars.dpoTableID + ' tbody tr');
            editRows.each((i) => {
                $("[id^='country-name-row-" + i + "']")[0].innerHTML = definitionsDpoObj[i]["CountryName"];
                $("[id^='country-code-row-" + i + "']")[0].innerHTML = definitionsDpoObj[i]["CountryCode"];
                $("[id^='country-default-row-" + i + "']")[0].value = definitionsDpoObj[i]["IsDefault"];
            });
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function closeJsonEditDpoModal() {
    jsonEditModalDpoOutputFieldId = "";
    $("#dpo-config-edit-modal").modal('hide');
}

let modalSoftUIVars;
let jsonEditModalSoftUIOutputFieldId;
let softUIModalCheck;

function showJsonEditSoftUIModal(modalTemplateEndpoint, elementName, outputFieldId, options) {

    softUIModalCheck = 1
    $.ajax({
        url: modalTemplateEndpoint,
        data: {
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            jsonEditModalSoftUIOutputFieldId = outputFieldId;
            $("#soft-ui-config-edit-modal-title").text(elementName);
            $("#soft-ui-config-edit-modal-body").html(clean);
            $("#soft-ui-config-edit-modal").show()
            $("#soft-ui-config-edit-modal").modal('show');
            $(".modal-backdrop").show();
            renderSoftUITable(options);
            let definitionsSoftUIJson =     $('#' + outputFieldId).val();
            let definitionsSoftUIObj = JSON.parse(definitionsSoftUIJson)
            let editRows = $(modalSoftUIVars.softUITableID + ' tbody tr');
            editRows.each((i) => {
                $("[id^='apm-row-" + i + "']")[0].value = definitionsSoftUIObj[i]["APM"];
                $("[id^='mcc-row-" + i + "']")[0].innerHTML = definitionsSoftUIObj[i]["MCC"];
            });
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

function closeJsonEditSoftUIModal() {
    jsonEditModalSoftUIOutputFieldId = "";
    $("#soft-ui-config-edit-modal").modal('hide');
}

function emptyToolTipAndCloseSoftUIModal() {
    $(modalSoftUIVars.dpoTableID + ' td[aria-describedby*="tooltip"]').hide();
    $("#soft-ui-config-edit-modal").modal('hide');
    softUIModalCheck = 0

}

function closeGratuityJsonEditModal() {
    $("#gratuity-edit-modal").modal('hide');
}

function showGratuityJsonEditModal(modalTemplateEndpoint, elementName, outputFieldId) {
    let fieldVal = document.getElementById(outputFieldId).value;
    $.ajax({
        url: modalTemplateEndpoint,
        data: {
            data: fieldVal,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: "POST",
        success: function (data) {
            hideWarning("gratuity");
            let clean = sanitizeHTML(data);
            $("#gratuity-edit-modal-title").text(elementName);
            $("#gratuity-edit-modal-body").html(clean);
            $("#gratuity-edit-modal").modal('show');
            if (fieldVal !== null || fieldVal !== "") {
                let jsonObject = JSON.parse(fieldVal);
                writeDataToModel(jsonObject)
            }
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    });
}

/**
 * Bind services data to loaded modal.
 * @param servicesJson The services data which comes as json string.
 */
function bindServiceData(servicesJson) {
    $('#error_midtid').hide()
    paymentServices = initiatePaymentServicesServices(servicesJson)
    //Make duplicate services array in the case of close the dialog it support to ignore the changes made.
    paymentServicesTemp = JSON.parse(JSON.stringify(paymentServices))
    const tableBody = $('#serviceTable tbody')
    //Bind the data to table
    paymentServicesTemp.forEach(function (service, index) {
        const newRow = $('<tr>');
        let cellMid = $('<td/>');
        let cellTid = $('<td/>');

        let midCallBack = function () {
            cellMid.children().remove();
            const midInput = $('<input>').attr('type', 'text').addClass('form-control float-right');
            midInput.on('blur', function () {
                cellMid.children().remove();
                $("<span/>").text(service.MID == "" ? "N/A" : service.MID).on('click', midCallBack).appendTo(cellMid);
            });
            const midError = $('<div>').addClass('invalid-feedback').text('The MID must be alphanumeric and have a length between 6 and 15 characters.');
            midInput.on('input', function () {
                $(this).val($(this).val().replace(/[^a-zA-Z0-9]/g, ''));
                const mid = $(this).val()
                if (regexMid.test(mid)) {
                    $(this).removeClass('is-invalid');
                    $(this).addClass('is-valid');
                    $('#error_midtid').hide()
                } else {
                    $(this).removeClass('is-valid');
                    $(this).addClass('is-invalid');
                }
                service.MID = mid
            });
            $(midInput).val(service.MID);
            cellMid.append(midInput, midError);
        };

        $("<span/>").text(service.MID == "" ? "N/A" : service.MID).on('click', midCallBack).appendTo(cellMid);

        let tidCallback = function () {
            cellTid.children().remove();
            const tidInput = $('<input>').attr('type', 'text').addClass('form-control float-right');
            tidInput.on('blur', function () {
                cellTid.children().remove();
                $("<span/>").text(service.TID > 0 ? service.TID : "N/A").on('click', tidCallback).appendTo(cellTid);
            });
            const tidError = $('<div>').addClass('invalid-feedback').text('The TID must be numeric, consist of 8 characters, and cannot be all zeros.')
            tidInput.on('input', function () {
                const tid = $(this).val()
                if (regexTid.test(tid)) {
                    $(this).removeClass('is-invalid');
                    $(this).addClass('is-valid');
                    $('#error_midtid').hide()
                } else {
                    $(this).removeClass('is-valid');
                    $(this).addClass('is-invalid');
                }
                service.TID = tid
            });
            $(tidInput).val(service.TID);
            cellTid.append(tidInput, tidError);
        };

        $("<span/>").text(service.TID > 0 ? service.TID : "N/A").on('click', tidCallback).appendTo(cellTid);

        // Append input boxes to the table row
        newRow.append($('<td>').append(service.Name));
        newRow.append(cellMid)
        newRow.append(cellTid)

        tableBody.append(newRow);
    })
    //Apply pagination table to service modal table

    $('#serviceTable').DataTable({ search: false, responsive: true });
    $('.dataTables_length').addClass('bs-select')
}

/**
 * Intialize if the paymentServices is already not defined else the function will return already defined paymentServices.
 * @param servicesJson The payment services as json string which received from back-end.
 * @returns {any|*[]} The payment services as array of objects.
 */
function initiatePaymentServicesServices(servicesJson) {
    if (paymentServices.length === 0) {
        let services = JSON.parse(servicesJson)
        //Assign already configured items data to the newly fetched services
        const configuredServicesJson = $('#' + serviceInputElementId).val();
        if (configuredServicesJson === "") {
            return services;
        }
        let configuredServices = []
        configuredServices = JSON.parse(configuredServicesJson);
        let servicesMap = {}
        services.forEach(function (newService, index1) {
            servicesMap[Math.round(newService.ServiceId)] = newService;
        })
        configuredServices.forEach(function (configService, index2) {
            let key = Math.round(configService.ServiceId);
            if (servicesMap.hasOwnProperty(key)) {
                servicesMap[key].MID = configService.MID
                servicesMap[key].TID = configService.TID
            }
        })
        return services;
    }
    return paymentServices;
}

/**
 * Load payment services configuration dialog(modal)
 * @param modalTemplateEndpoint The url which can request modal template
 * @param siteId The site id which services field contain.
 * @param tid The terminal id
 * @param serviceInputId The id of the hidden input element which stored configured data as json
 * @param elementName The dialog title.
 */
function showPaymentServicesEditModal(modalTemplateEndpoint, siteId, tid, serviceInputId, elementName) {
    if (serviceInputId !== serviceInputElementId) {
        resetPaymentData()
    }
    serviceInputElementId = serviceInputId
    $.ajax({
        url: modalTemplateEndpoint,
        data: {
            siteId: siteId,
            tid: tid,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: "POST",
        success: function (d) {
            var cleanData = sanitizeHTML(d) //remove malicious html elements
            $("#services-edit-modal-title").text(elementName)
            $("#services-edit-modal-body").html(cleanData)
            $("#services-edit-modal").modal('show')

            var serviceJson = $('#serviceJson').val()
            bindServiceData(serviceJson)
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });
}

/**
 * Save the payment services locally in the hidden input filed which named to save the services data automatically.
 */
function savePaymentServices() {
    //Apply the configured payment services from temperately array to permanent array
    paymentServices = paymentServicesTemp
    paymentServicesTemp = []
    let nonEmptyServices = []
    $('#error_midtid').hide()
    const seenMids = new Set();
    const seenTids = new Set();
    for (let i = 0; i < paymentServices.length; i++) {
        let service = paymentServices[i]
        if (service.TID !== "" && service.TID !== "N/A" && service.MID !== "" && service.MID !== "N/A") {
            if ((!regexTid.test(service.TID + '')) || !regexMid.test(service.MID)) {
                showPaymentServiceError("Invalid MID/TID detected. Please review your input for accuracy")
                paymentServicesTemp = paymentServices
                return
            }

            if (seenMids.has(service.MID)) {
                showPaymentServiceError("The MID " + service.MID + " duplicated.")
                paymentServicesTemp = paymentServices
                return
            } else {
                seenMids.add(service.MID)
            }

            if (seenTids.has(service.TID)) {
                showPaymentServiceError("The TID " + service.TID + " duplicated.")
                paymentServicesTemp = paymentServices
                return
            } else {
                seenTids.add(service.TID)
            }
            nonEmptyServices.push(service)
        }
    }
    $('#' + serviceInputElementId).val(JSON.stringify(nonEmptyServices))
    $("#services-edit-modal-close").click()
}

function closeServicesModal() {
    $("#services-edit-modal").modal('hide');
    configModalCheck = 0
}

function showPaymentServiceError(error) {
    $('#error_midtid').text(error)
    $('#error_midtid').show()
}

/**
 * Cancel the payment service configurations and close the dialog.
 */
function cancelServiceConfig() {
    paymentServicesTemp = []
    $("#services-edit-modal-close").click()
}

function resetPaymentData() {
    paymentServicesTemp = []
    paymentServices = []
    serviceInputElementId = null
}

function writeDataToModel(fieldVal) {
    $.each(fieldVal, function (i, item) {
        $("#" + i).val(item)
        if (i === "tipProType"){
            changeTIPConfiguration(item)
        }
    });
}

function validNumber(preDefineTip) {
    const regex = /(^[0-9]\d{0,4})+$|^$|^\s$/gm; //Ensure only positive integer values from between 0-99999 (AED 999.99 on PED)
    return regex.test(preDefineTip);
}

function ValidPercentage(preDefineTipPercentage){
    const regexPercentage = new RegExp('^(0*100{1,1}\\.?((?<=\\.)0*)?%?$)|(^0*\\d{0,2}\\.?((?<=\\.)\\d*)?%?)$', 'mgi')
    return regexPercentage.test(preDefineTipPercentage);

}


function changeTIPConfiguration(type){
    if (type.value === "1" || type === "1"){
        document.getElementById("preDefineTipPercentage").disabled = true;
        document.getElementById("preDefineTip").disabled = false;
    } else if (type.value === "2" || type === "2"){
        document.getElementById("preDefineTip").disabled = true;
        document.getElementById("preDefineTipPercentage").disabled = false;
    }
}

function updateGratuity() {
    let object = {};
    hideWarning("gratuity");
    let tipProType = document.getElementById('tipProType').value;
    let textCharacterLimit = document.getElementById('textCharacterLimit').value;
    let english = document.getElementById('english').value;
    let french = document.getElementById('french').value;
    let arabic = document.getElementById('arabic').value;
    object["tipProType"] = tipProType;
    object["textCharacterLimit"] = textCharacterLimit;
    object["english"] = english;
    object["french"] = french;
    object["arabic"] = arabic

    if (tipProType === "") {
        displayWarningMessage('["tipProType:Please Enter Select Tip Pro Type"]', "gratuity")
        return false;
    }

    if (textCharacterLimit == null || textCharacterLimit === "" || textCharacterLimit === 0 || textCharacterLimit > 100) {
        displayWarningMessage('["TextCharacterLimit:Please Enter valid Text Character Limit Between 1-100 only"]', "gratuity")
        return false;
    }

    if (english === "") {
        displayWarningMessage('["English:Please Enter English Message"]', "gratuity")
        return false;
    }

    if (english.length > textCharacterLimit) {
        displayWarningMessage('["English:Please ensure that the text doesn\'t exceed the text character limit value"]', "gratuity")
        return false;
    }

    if (french !== "" && french.length > textCharacterLimit) {
        displayWarningMessage('["French:Please ensure that the text doesn\'t exceed the text character limit value"]', "gratuity")
        return false;
    }

    if (arabic !== "" && arabic.length > textCharacterLimit) {
        displayWarningMessage('["Arabic:Please ensure that the text doesn\'t exceed the text character limit value"]', "gratuity")
        return false;
    }

    for (let i = 1; i <= 10; i++) {
        let preDefineTip = document.getElementById('preDefineTip' + i).value;
        if (validNumber(preDefineTip) === false) {
            displayWarningMessage(`["PreDefineTip${i}: Please enter a valid value up to 99999"]"]`, "gratuity");
            return false;
        }
        if (preDefineTip !== "") {
            object[`preDefineTip${i}`] = Number(preDefineTip)
        }
    }

    for (let j = 1; j <= 10; j++) {
        let preDefineTipPercentage = document.getElementById('preDefineTipPercentage' + j).value;
        if (ValidPercentage(preDefineTipPercentage) === false) {
            displayWarningMessage(`["PreDefineTipPercentage${j}: Please enter a valid value up to 100"]"]`, "gratuity");
            return false;
        }
        if (preDefineTipPercentage !== "") {
            object[`preDefineTipPercentage${j}`] = Number(preDefineTipPercentage)
        }
    }

    let json = JSON.stringify(object);
    $("#modules-gratuityConfigs").val(json)
    closeGratuityJsonEditModal()
}

function bindChooseFile(elementId, directory) {
    $('[data-type="choose-file-name"]').click(function () {
        let fileName = $(this).attr('data-filename')
        let directory = $(this).attr('data-filetype')
        $('#' + elementId).val(fileName)
        $('#chooseFileModal').modal('hide')

        const data = {
            FileName: fileName,
            Directory: directory,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        }

        $.ajax({
            url: 'getFile',
            data: data,
            method: 'POST',
            success: function (d) {
                let clean = sanitizeHTML(d);
                $('#' + elementId + '-image').attr('src', clean.replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/'/g, '&#39;').replace(/</g, '&lt;').replace(/>/g, '&gt;'))
            },
            error: function (d) {
                let clean = DOMPurify.sanitize(d.responseText)
                displayWarningMessage(clean)
            }
        })
    })
}

function SaveTidOverride(tidId) {
    $(".xdsoft_datetimepicker").remove();
    let dataObj = new FormData($("#saveTidProfile_" + tidId)[0]);
    dataObj.append('pageSize', $("#pageSize")[0].value)
    dataObj.append('pageNumber', $("#pageNumber")[0].value)
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value)
    resetPaymentData()

    $.ajax({
        data: dataObj,
        method: "POST",
        url: "saveTidProfile",
        processData: false,
        contentType: false,

        success: function (data) {
            let clean = sanitizeHTML(data);
            $('#tids').html(clean)
            hideWarning();
            $('[name^=multi]').multiselect();
            bindSubmit();
            bindShowTIDDetails();
            bindShowFraudOverride();
            bindShowUserOverride();
            bindOverrideField();
            bindToggle();
            $('[data-button="toggle"]').parents().next('.hide').toggle();
        },
        error: function (data) {
            //the data is a json array of messages
            if (data.getResponseHeader("content-type") === "application/json") {
                let validationMessages = JSON.parse(data.responseText).join('<br/>');
                displayWarningMessage(validationMessages)
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    })
}

function DeleteTidOverride(tidId) {
    $(".xdsoft_datetimepicker").remove();
    let dataObj = new FormData($("#saveTidProfile_" + tidId)[0]);
    dataObj.append('pageSize', $("#pageSize")[0].value)
    dataObj.append('pageNumber', $("#pageNumber")[0].value)
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value)

    resetPaymentData()

    $.ajax({
        data: dataObj,
        method: "POST",
        url: "deleteTidProfile",
        processData: false,
        contentType: false,

        success: function (data) {
            let clean = sanitizeHTML(data);
            $('#tids').html(clean)
            hideWarning();
            $('[id^=multiselect]').multiselect();
            bindSubmit();
            bindOverrideField();
            bindToggle();
            $('[data-button="toggle"]').parents().next('.hide').toggle();
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function AddDuplicatedTidOverride(tidId) {
    resetPaymentData()
    $(".xdsoft_datetimepicker").remove();
    let dataObj = $("#saveTidProfile_" + tidId).serializeArray()
    dataObj.push({ name: "TidIndex", value: 0 })
    dataObj.push({ name: "Site", value: $("#siteID")[0].value })
    dataObj.push({ name: "SiteProfileId", value: $("#profileID")[0].value })
    CreateNewTidWithDuplicateOverride(dataObj)
    $("#tid-override-" + tidId).toggle()
}

function ApplyDuplicateNewTidOverride(tidId) {
    resetPaymentData()
    $(".xdsoft_datetimepicker").remove();
    let dataObj = new FormData($("#addTidProfile_" + tidId)[0]);
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value)

    let dObj = $("#addTidProfile_" + tidId).serializeArray();
    dObj.push({ name: "TidIndex", value: 0 });
    dObj.push({ name: "Site", value: $("#siteID")[0].value });
    dObj.push({ name: "SiteProfileId", value: $("#profileID")[0].value })

    $.ajax({
        url: "addTID",
        data: dataObj,
        method: "POST",
        processData: false,
        contentType: false,
        success: function (d) {
            let clean = sanitizeHTML(d);
            $('#tids').html(clean);
            hideWarning();
            $('[id^=multiselect]').multiselect();
            bindSubmit();
            bindOverrideField();
            bindShowFraudOverride();
            bindShowUserOverride();
            bindShowTIDDetails();
            bindToggle();
            $('[data-button="toggle"]').parents().next('.hide').toggle();
            CreateNewTidWithDuplicateOverride(dObj)
        },
        error: function (data) {
            //the data is a json array of messages
            if (data.getResponseHeader("content-type") === "application/json") {
                let validationMessages = JSON.parse(data.responseText).join('<br/>');
                displayWarningMessage("Error saving TID: " + validationMessages);
            } else {
                displayWarningMessage("Error saving TID: " + data.responseText);

            }
        }
    });
}

function CreateNewTidWithDuplicateOverride(dObj) {
    resetPaymentData()
    $.ajax({
        url: "addNewDuplicatedTidOverride",
        data: dObj,
        method: "POST",
        success: function (d) {
            let clean = sanitizeTableHTML(d);
            $("#tid_table").prepend(clean)
            $('[name^=multi]').multiselect();
        },
        error: function (data) {
            //the data is a json array of messages
            if (data.getResponseHeader("content-type") === "application/json") {
                let validationMessages = JSON.parse(data.responseText).join('<br/>');
                displayWarningMessage("Error saving TID: " + validationMessages)
            } else {
                displayWarningMessage("Error saving TID: " + data.responseText)
            }
        }
    });
}

function CancelAddTid(tidId) {
    $(".xdsoft_datetimepicker").remove();
    let row = $("#addTidProfile_" + tidId)[0];
    document.getElementById("tid_table").deleteRow(row);
}

function checkIsDirty(tidId) {
    $("#duplicate-" + tidId).attr('disabled', 'disabled');
}


function eodAutoTypeChange() {
    var e = document.getElementById("endOfDayAutoType");
    var text = e.options[e.selectedIndex].text;
    var t = document.getElementById("endOfDay-time");

    if (text === "Time Range") {
        $("#endOfDayAutoTime").html(timeElement(t, "TimeRange"));
    } else {
        $("#endOfDayAutoTime").html(timeElement(t, "SpecificTime"));
    }
}

function timeElement(t, timeType) {
    if (timeType == "TimeRange") {
        var from = `From : <input id="endOfDay-time" type="text" name="` + t.getAttribute("name") + `" placeholder="From"`
        var to = ` To : <input id="endOfDay-time-2" type="text" name="` + t.getAttribute("name") + `***" placeholder="To"`

        if (t.getAttribute("pattern") != null) {
            from += `pattern="` + t.getAttribute("pattern") + `"`
            to += `pattern="` + t.getAttribute("pattern") + `"`
        }

        if (t.getAttribute("oninvalid") != null) {
            from += `oninvalid="` + t.getAttribute("oninvalid") + `"`
            to += `oninvalid="` + t.getAttribute("oninvalid") + `"`
        }

        if (t.getAttribute("oninput") != null) {
            from += `oninput="` + t.getAttribute("oninput") + `"`
            to += `oninput="` + t.getAttribute("oninput") + `"`
        }

        from += `/>`
        to += `/>`
        return (from + to)
    }
    else {
        var h = `<input id="endOfDay-time" type="text" class="form-control" name="` + t.getAttribute("name") + `" placeholder="Time"`

        if (t.getAttribute("pattern") != null) {
            h += `pattern="` + t.getAttribute("pattern") + `"`
        }

        if (t.getAttribute("oninvalid") != null) {
            h += `oninvalid="` + t.getAttribute("oninvalid") + `"`
        }

        if (t.getAttribute("oninput") != null) {
            h += `oninput="` + t.getAttribute("oninput") + `"`
        }

        h += `/>`
        return h
    }
}

//# sourceURL=profileMaintenance.js