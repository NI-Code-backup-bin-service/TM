$(document).ready(() => {
    showBulkChangeApprovalApprove();
})


function showBulkChangeApprovalApprove() {
    hideWarning();

    $.ajax({
        url: 'bulkChangeApproval/unapproved',
        success: (data) => {
            let clean = sanitizeHTML(data);
            $('#bulkChangeApprovalContainerApprove').show();
            $('#bulkChangeApprovalContainerApprove').html(clean);
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}

function showBulkChangeApprovalHistory() {
    hideWarning();
    const data = {
        csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    }
    $.ajax({
        url: 'bulkChangeApproval/history',
        data: data,
        success: (data) => {
            let clean = sanitizeHTML(data);
            $('#bulkChangeApprovalContainerHistory').show();
            $('#bulkChangeApprovalContainerHistory').html(clean);
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}

function downloadFile(fileName, fileType, report) {
    hideWarning();
    if (report) {
        fileName = "report_" + fileName
    }

    const data = {
        FileName: fileName,
        FileType: fileType,
        csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    }

    $.ajax({
        url: 'downloadFile',
        data: data,
        method: 'POST',
        success: (d, status, req) => {
            const link = document.createElement('a')
            link.href = encodeURI(window.URL.createObjectURL(new Blob([d], {
                type: 'text/plain'
            })));
            if (report) {
                // filename will have date as well, So removing the date from the file name (yyyymmddhhmmss_)
                link.download = "report_" + fileName.substring(22);
            } else {
                // filename will have date as well, So removing the date from the file name (yyyymmddhhmmss_)
                link.download = fileName.substring(15);
            }
            link.click()
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}

function approve(fileName, fileType, changeType) {
    hideWarning();

    const data = {
        FileName: fileName,
        FileType: fileType,
        ChangeType: changeType,
        csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    }

    $.ajax({
        url: 'bulkChangeApproval/approve',
        data: data,
        method: 'POST',
        success: () => {
            showBulkChangeApprovalApprove()
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}

function discard(fileName, fileType, changeType) {
    hideWarning();

    const data = {
        FileName: fileName,
        FileType: fileType,
        ChangeType: changeType,
        csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    }

    $.ajax({
        url: 'bulkChangeApproval/discard',
        data: data,
        method: 'POST',
        success: () => {
            showBulkChangeApprovalApprove()
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}