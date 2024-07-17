$(document).ready(() => {
    // Hide the spinner on page load
    $('#site-upload-spinner').hide();
    setupUploadListener();
})


function setupUploadListener() {
    $('#terminal-flagging-import-upload').submit(function(e) {
        // When the form is submitted immediately show the spinner
        hideWarning();
        $('#site-upload-spinner').show();
        $('#terminal-flagging-upload-results').hide();
        e.preventDefault();

        const formData = new FormData(document.querySelector('#terminal-flagging-import-upload'));
        $.ajax({
            type: "POST",
            data: formData,
            url: "/terminalFlagging/upload",
            processData: false,
            contentType: false,
            success: function(data) {
                let clean = sanitizeHTML(data);
                // Hide the spinner when the data is returned
                $('#site-upload-spinner').hide();
                $('#terminal-flagging-upload-results').show();
                $('#terminal-flagging-upload-results').html(clean);
                $('#terminal-flagging-import-upload-file').val(null);
                $('#terminal-flagging-upload-results').delay(2000).fadeOut('slow');
            },

            error: function(data) {
                // Hide the spinner if an error is returned
                $('#site-upload-spinner').hide();
                displayWarningMessage(data.responseText);
            }
        });
    });
}
