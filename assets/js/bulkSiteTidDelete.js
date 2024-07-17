function bulkDeleteTid() {
    $("#tid-delete-results").hide()
    hideWarning();
    let formData = new FormData($('#tid-delete-upload')[0]);
    formData.append("DeleteType", 'TidDelete')
    $.ajax({
        data: formData,
        type: "POST",
        url: "bulkDelete",
        processData: false,
        contentType: false,
        success: function (data) {
            let clean = sanitizeHTML(data)
            $('#tid-delete-results').hide()
            $("#tid-delete-results").html(clean);
            $("#tid-delete-results").show();
            $('#tid-delete-upload-file').val(null);
        },
        error: function (data) {
            $('#tid-delete-results').hide()
            displayWarningMessage(data.responseText)
        }
    })
}

function bulkDeleteSite() {
    let formData = new FormData($('#site-delete-upload')[0]);
    formData.append("DeleteType", 'SiteDelete')
    $.ajax({
        data: formData,
        type: "POST",
        url: "bulkDelete",
        processData: false,
        contentType: false,
        success: function (data) {
            let clean = sanitizeHTML(data)
            $('#site-delete-results').hide()
            $("#site-delete-results").html(clean);
            $("#site-delete-results").show();
            $('#site-delete-upload-file').val(null);
        },
        error: function (data) {
            $('#site-delete-results').hide()
            displayWarningMessage(data.responseText)
        }
    })

}