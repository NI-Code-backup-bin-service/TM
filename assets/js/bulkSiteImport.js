$(document).ready(() => {
    // Hide the spinner on page load
    $('#site-upload-spinner').hide();
    setupUploadListener();
})

// String constants
let sitesCommittedSuccess = "All sites successfully created and committed to database";
let commitConfirmationTitle = "Commit Bulk Site Import?";
let commitConfirmationBody = "Are you sure you want to commit all of the validated sites?";
let cancelConfirmationTitle = "Cancel Bulk Site Import?";
let cancelConfirmationBody = "Are you sure you want to cancel? This will reset all completed validation.";

let modalShown = false;
let commitMade = false;

function setupUploadListener() {
    $('#site-import-upload').submit( function(e) {
        // When the form is submitted immediately show the spinner
        $('#site-upload-spinner').show();
        hideWarning();
        $("#site-upload-results").hide();
        e.preventDefault();
        const formData = new FormData(document.querySelector('#site-import-upload'));

        $.ajax({
            data: formData,
            type: "POST",
            url: "uploadSites",
            processData: false,
            contentType: false,

            success: function (data){
                let clean = sanitizeHTML(data);

                // Hide the spinner when the data is returned
                $('#site-upload-spinner').hide();
                $("#site-upload-results").show();
                $("#site-upload-results").html(clean);
                $("#site-upload-results-table tr").click(function () {
                    $(this).toggleClass("cap-height");
                });
            },

            error: function(data){
                // Hide the spinner if an error is returned
                $('#site-upload-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        })
    })
}

function confirmSiteImportDialog(title, msg, onAccept, onDecline){
    modalShown = true;
    var p = showConfirmDialog(title,msg);
    p.done( function(confirmed){
        if(confirmed){
            onAccept();
        } else {
            onDecline();
        }
    })

}

// Used to capture instances when user closes modal by clicking outside of the modal
document.addEventListener("click", (evt) => {
    if (modalShown) {
        // Defines the modal dialogue box
        const modal = document.getElementById("confirmModalContent");
        const buttonClicked = document.getElementById("commit-site-imports");
        let targetElement = evt.target;
        do {
            if (targetElement === modal) {
                // If the click takes place inside the dialogue box, do nothing
                return;
            } else if (targetElement === buttonClicked) {
                // If the click is on the commit button, make sure the buttons are set to disabled
                toggleSiteValidationButtons(true);
                return;
            }
            targetElement = targetElement.parentNode;
        } while (targetElement);

        // If the click does not fall within the modal or is not on the commit button then it will close
        // When the modal closes we need to re-enable the buttons
        toggleSiteValidationButtons(false);
    }
})

// Enables/disables the site import validation buttons
function toggleSiteValidationButtons(disable) {
    $('#commit-site-imports').prop('disabled', disable);
    $('#cancel-site-imports').prop('disabled', disable);
}

function commitBulkSiteImport() {
    toggleSiteValidationButtons(true);
    // Show a confirmation dialogue box as this process can potentially store thousands of new sites
    confirmSiteImportDialog(commitConfirmationTitle, commitConfirmationBody, function() {
        // When upload is committed immediately show the spinner
        $('#site-upload-spinner').show();
        // In here we want to take the new sites and send them back to be stored in the DB
        let commitCall = $.ajax({
            data: {
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
            },
            type: "POST",
            url: "/commitBulkSiteUpload",
            beforeSend: function () {
                // If the user repeatedly opens and closes the commit modal Ajax will queue up multiple calls and then if the
                // user clicks OK to commit ALL of the queued calls will be sent. This function will check to see if a commit
                // call has already been made, and if it has it will abort additional calls.
                if (commitMade) {
                    commitCall.abort();
                }
                commitMade = true;
            },
            success: function (data){
                // Hide the spinner when the data is returned
                $('#site-upload-spinner').hide();
                commitMade = false;
                // Re-enable the buttons to prevent any issues on next visit
                toggleSiteValidationButtons(false);
                clearSiteInputs();
                displayWarningMessage(sitesCommittedSuccess)
            },

            error: function(data){
                // Hide the spinner if an error is returned
                $('#site-upload-spinner').hide();
                commitMade = false;
                toggleSiteValidationButtons(false);
                displayWarningMessage(data.responseText)
            }
        })
    }, function() {
        toggleSiteValidationButtons(false)
        modalShown = false;
    })
}

function clearSiteInputs() {
    // Stop displaying the validation results
    $("#site-upload-results").html(null);
    // Clear the template MID
    $('#site-import-template-MID').val("");
    // Clear the file selector
    $('#site-import-upload-file').val(null);
}

function cancelBulkSiteImport() {
    // Show a confirmation dialogue box as this process can potentially store thousands of new sites
    confirmDialog(cancelConfirmationTitle, cancelConfirmationBody, function() {
        clearSiteInputs();
    })
}

