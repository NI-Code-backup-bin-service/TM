$(document).ready(() => {
    $('#site-update-spinner').hide();
})

function updateSites() {
        $('#site-update-results').hide()
        $('#site-update-spinner').show();
        hideWarning();
        $('#site-update-spinner').hide();
        const formData = new FormData(document.querySelector('#site-update-upload'));
        $.ajax({
            data: formData,
            type: "POST",
            url: "/updateSites",
            processData: false,
            contentType: false,
            success: function (data) {
                let clean = sanitizeHTML(data)
                $('#site-update-spinner').hide();
                $("#site-update-results").html(clean);
                $('#site-update-results').show()
                $('#site-update-upload-file').val(null);
            },

            error: function (data) {
                $('#site-update-results').hide()
                $('#site-update-spinner').hide();
                displayWarningMessage(data.responseText)
            }
        });
}
