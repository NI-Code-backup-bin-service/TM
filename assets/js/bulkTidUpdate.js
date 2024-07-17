$(document).ready(() => {
    // Hide the spinner on page load
    $('#tid-update-spinner').hide();
})

function updateTids() {
    $("#tid-update-results").hide()
    $('#tid-update-spinner').show();
    hideWarning();
    $('#tid-update-spinner').hide();
    const formData = new FormData(document.querySelector('#tid-update-upload'));
    $.ajax({
        data: formData,
        type: "POST",
        url: "updateTids",
        processData: false,
        contentType: false,
        success: function (data) {
            let clean = sanitizeHTML(data)
            $('#tid-update-results').hide()
            $('#tid-update-spinner').hide();
            $("#tid-update-results").html(clean);
            $("#tid-update-results").show();
            $('#tid-update-upload-file').val(null);
        },
        error: function (data) {
            $('#tid-update-results').hide()
            $('#tid-update-spinner').hide();
            displayWarningMessage(data.responseText)
        }
    })
}