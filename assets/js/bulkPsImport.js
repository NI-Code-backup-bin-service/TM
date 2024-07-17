$(document).ready(() => {
    $('#ps-upload-spinner').hide();
    hideResults();
    setupPsGroupImport();
    setupPsTidImport();
})

function clearInputs() {
    $('#ps-groups-import-upload-file').val(null);
    $('#ps-tid-import-upload-file').val(null);
}

function hideResults() {
    $("#ps-upload-results").hide();
    $("#ps-tid-upload-results").hide();
}
function setupPsGroupImport() {
    $('#ps-groups-import-upload').submit(function (event) {
        hideResults();
        $('#ps-upload-spinner').show();
        event.preventDefault();
        const formData = new FormData(document.querySelector('#ps-groups-import-upload'));

        $.ajax({
            data: formData,
            type: "POST",
            url: "uploadPaymentServicesGroups",
            processData: false,
            contentType: false,
            success: function (data) {
                let clean = sanitizeHTML(data)
                $('#ps-upload-spinner').hide();
                $("#ps-upload-results").html(clean);
                $('#ps-upload-results').show()
            },

            error: function (data) {
                $('#ps-upload-results').hide();
                $('#ps-upload-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function commitPsGroupImport() {
    let token = $('input[name=csrfmiddlewaretoken]')[0].value
    $.ajax({
        type: "POST",
        data: {
            csrfmiddlewaretoken: token
        },
        url: "commitPaymentServicesGroups",

        success: function (data) {
            data = JSON.parse(data)
            $('#ps-upload-spinner').hide();
            hideResults();
            clearInputs();

            let message = "no new rows were imported into the database"
            if (data && data.importedRows) {
                message = `successfully imported ${data.importedRows} rows into the database`
            }
            displayWarningMessage(message)
        },

        error: function (data) {
            $('#ps-upload-spinner').hide();
            hideResults();
            clearInputs();
            displayWarningMessage(data.responseText)
        }
    })
}

function setupPsTidImport() {
    $('#ps-tid-import-upload').submit(function (event) {
        hideResults();
        $('#ps-upload-spinner').show();
        event.preventDefault();
        const formData = new FormData(document.querySelector('#ps-tid-import-upload'));

        $.ajax({
            data: formData,
            type: "POST",
            url: "uploadPaymentServicesTids",
            processData: false,
            contentType: false,

            success: function (data) {
                let clean = sanitizeHTML(data)
                $('#ps-upload-spinner').hide();
                $("#ps-tid-upload-results").html(clean);
                $('#ps-tid-upload-results').show()
            },

            error: function (data) {
                $('#ps-tid-upload-results').hide()
                $('#ps-upload-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function commitPsTidImport() {
    let token = $('input[name=csrfmiddlewaretoken]')[0].value
    $.ajax({
        type: "POST",
        data: {
            csrfmiddlewaretoken: token
        },
        url: "commitPaymentServicesTids",

        success: function (data) {
            data = JSON.parse(data)
            $('#ps-upload-spinner').hide();
            hideResults();
            clearInputs();

            let message = "no new rows were imported into the database"
            if (data && data.importedRows) {
                message = `successfully imported configuration for ${data.importedRows} terminals`
            }
            displayWarningMessage(message)
        },

        error: function (data) {
            $('#ps-upload-spinner').hide();
            hideResults();
            clearInputs();
            displayWarningMessage(data.responseText)
        }
    })
}

function cancelPsImport() {
    confirmDialog("Cancel Bulk Import?", "Are you sure you want to cancel? This will reset all completed validation.", function () {
        clearInputs();
        hideResults();
    })
}
