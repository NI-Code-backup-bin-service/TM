$(document).ready(() => {
    // Hide the spinner on page load
    $('#tid-upload-spinner').hide();
    setupTidUploadListener();
})

// String constants
let tidsCommittedSuccess = "All TIDs successfully created and committed to database";
let tidCommitConfirmationTitle = "Commit Bulk TID Import?";
let tidCommitConfirmationBody = "Are you sure you want to commit all of the validated TIDs to their respective sites?";
let tidCancelConfirmationTitle = "Cancel Bulk TID Import?";
let tidCancelConfirmationBody = "Are you sure you want to cancel? This will reset all completed validation.";

function setupTidUploadListener() {
    $('#tid-import-upload').submit( function(e) {
        // When the form is submitted immediately show the spinner
        $('#tid-upload-spinner').show();
        e.preventDefault();
        const formData = new FormData(document.querySelector('#tid-import-upload'));
        $.ajax({
            data: formData,
            type: "POST",
            url: "uploadTids",
            processData: false,
            contentType: false,

            success: function (data){
                let clean = sanitizeHTML(data)
                // Hide the spinner when the data is returned
                hideWarning()
                $("#tid-upload-results").hide()
                $('#tid-upload-spinner').hide();
                $("#tid-upload-results").html(clean);
                $("#tid-upload-results").show()
                $("#tid-upload-results-table tr").click(function () {
                    $(this).toggleClass("cap-height");
                });
            },

            error: function(data){
                // Hide the spinner if an error is returned
                $("#tid-upload-results").hide()
                $('#tid-upload-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function commitBulkTidImport() {
    // Show a confirmation dialogue box as this process can potentially store thousands of new sites
    confirmDialog(tidCommitConfirmationTitle, tidCommitConfirmationBody, function() {
        // When upload is committed immediately show the spinner
        $('#tid-upload-spinner').show();
        // In here we want to take the new sites and send them back to be stored in the DB
        $.ajax({
            data: {
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
            },
            type: "POST",
            url: "/commitBulkTidUpload",
            success: function (data){
                // Hide the spinner when the data is returned
                $('#tid-upload-spinner').hide();
                clearInputs();
                displayWarningMessage(tidsCommittedSuccess)
            },

            error: function(data){
                // Hide the spinner if an error is returned
                $('#tid-upload-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function clearInputs() {
    // Stop displaying the validation results
    $("#tid-upload-results").html(null);
    // Clear the file selector
    $('#tid-import-upload-file').val(null);
}

function cancelBulkTidImport() {
    // Show a confirmation dialogue box as this process can potentially store thousands of new sites
    confirmDialog(tidCancelConfirmationTitle, tidCancelConfirmationBody, function() {
        clearInputs();
    })
}