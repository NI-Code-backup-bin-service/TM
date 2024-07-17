$(document).ready(() => {
    $('#upload-file').submit( function(e) {
        e.preventDefault();
        var dataObj = new FormData(document.querySelector("#upload-file"))
        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            processData: false,
            contentType: false,
            success: function(data){
                window.location.href = "/search"
            },
            error: function(data){
                displayWarningMessage(data.responseText)
            }
        })
    })

    $('#upload-transactions').submit( function(e) {
        e.preventDefault();
        var dataObj = new FormData(document.querySelector("#upload-transactions"))
        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            processData: false,
            contentType: false,
            success: function(data){
                const exportTable = ce(ExportTable, {
                    id: 'export-table',
                    key:"txn-result-table",
                    Columns: ["Offline file", "Success", "Result"],
                    Rows: data.Results,
                    Filename: "Transactions"
                }, null);
                ReactDOM.render(exportTable, document.getElementById("txn-results"));
                $("#txnResultsModal").modal('show');
            },
            error: function(data){
                displayWarningMessage(data.responseText)
            }
        })
    })

    $('#softui-file-upload').submit( function(e) {
      e.preventDefault();
      const formData = new FormData(document.querySelector('#softui-file-upload'));
      $.ajax({
          data: formData,
          type: "POST",
          url: "softUiFileUpload",
          processData: false,
          contentType: false,
          success: function(data){
              window.location.href = "/search"
          },
          error: function(data){
              displayWarningMessage(data.responseText)
          }
      })
    })
})

function openDeleteFile (){
    var postData = {
        ButtonText: "Delete",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "getFileList",
        data:postData,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d);
            $("#choose-file-body").html(clean)
            bindDeleteFile()
            $("#chooseFileModal").modal('show')
        },
        error: function(d){
            displayWarningMessage(d.responseText)
        }
    });;
}

function bindDeleteFile(){
    $("[data-type='choose-file-name']").click(function () {
        var fileName = $(this).attr('data-filename')
        confirmDialog("Delete Site?",'Are you sure you want to delete '+fileName+'?', function(){
            var data = {
                fileName: fileName,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            }
            $.ajax({
                url: "deleteFile",
                data: data,
                method: "POST",
                success: function (d) {
                    $("#chooseFileModal").modal('hide')
                },
                error: function(data){
                    displayWarningMessage(data.responseText)
                }
            });;
        })
    })
}