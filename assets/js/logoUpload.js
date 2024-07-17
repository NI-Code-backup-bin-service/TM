$(document).ready(function () {
    var postData = {
        ButtonText: "Delete",
        type: "mnoLogo",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "getFileList?type=mnoLogo",
        data: postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindDeleteMnoLogo()
            $("#chooseFileDialog").modal('show')
        },
        error: function (d) {
            displayWarningMessage(d.responseText)
        }
    });

    $.ajax({
        url: "getDpoMomoFieldsData",
        method: "GET",
        success: function (result) {
            let resp = result.split("|").sort()
            let options = '<option id="mno-logo-select" value="select">--Select MNO--</option>'
            resp.forEach(function (value, key) {
                options += '<option id="mno-logo-' + value + '" value="' + value + '">' + value + '</option>'
            })
            $('#mno-logo').html(sanitizeHTML(options));
        },
        error: function (err) {
            displayWarningMessage(err.responseText)
        }
    });

    // below function used for dpoMomo logo upload
    $('#mno-logo-upload').submit(function (e) {
        e.preventDefault();
        var dataObj = new FormData(document.querySelector("#mno-logo-upload"))
        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            processData: false,
            contentType: false,
            success: function (data) {
                window.location.reload();
            },
            error: function (data) {
                displayWarningMessage(data.responseText)
            }
        })
    })
});

function bindDeleteMnoLogo() {
    $("[data-type='choose-file-name']").click(function () {
        var fileName = $(this).attr('data-filename')
        confirmDialog("Delete Site?", 'Are you sure you want to delete ' + fileName + '?', function () {
            var data = {
                fileName: fileName,
                directory: "mnoLogo",
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            }
            $.ajax({
                url: "deleteFile",
                data: data,
                method: "POST",
                success: function (d) {
                    window.location.reload();
                    $("#chooseFileModal").modal('hide')
                },
                error: function (data) {
                    displayWarningMessage(data.responseText)
                }
            });
        })
    })
}