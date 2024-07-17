
document.addEventListener('DOMContentLoaded', function() {
    $('#formContent').hide()
    $('#spinner').hide()
}, false);

// String constants
let sitesCommittedSuccess = "Data Successfully committed to database";
let commitConfirmationTitle = "Commit Data?";
let commitConfirmationBody = "Are you sure you want to commit all of the data?";
let cancelConfirmationTitle = "Cancel Submission?";

let modalShown = false;
let commitMade = false;

function confirmSubmitDialog(title, msg, onAccept, onDecline){
    modalShown = true;
    var p = showConfirmDialog(title,msg)
    p.done( function(confirmed){
        if(confirmed){
            onAccept();
        } else {
            onDecline();
        }
    })

}
function exportContactPDF() {
    var doc = new jsPDF({
        unit: 'pt',
        format: [595, 842]
    })
        html2canvas(document.getElementById("pageContentForm"), {
            useCORS: true,
            background: "#ffffff",
    }).then( function(canvas) {
            var myImage = canvas.toDataURL("image/jpeg,1.0");
            doc.addImage(myImage, 'png', 20, 10, 550, 550 )
            doc.save('contactSheet.pdf')})

}
function openForm(){
    $('#pageContent').hide()
    $('#formContent').show()
}
function closeForm(){
    $('#pageContent').show()
    $('#formContent').hide()
}

function submitForm(){

    const token = $("input[name=csrfmiddlewaretoken]")[0].value
    confirmSubmitDialog(commitConfirmationTitle, commitConfirmationBody, function(){
        $('#spinner').show();
        const formData = new Object();
        $("#formData input[type=text]").each(function(i, data) {
            var id = data.id;
            formData[id] = data.value
            console.log(data.id, data.value)
        });
        $("#formData textarea").each(function(i, data) {
            var id = data.id;
            formData[id] = data.value
            console.log(data.id, data.value)
        });
        // In here we want to take the new Contact us form  and send them back to be stored in the DB
        let commitCall = $.ajax({
            url: "/submitContactUsForm",
            data: {csrfmiddlewaretoken:token, form:formData},
            type: "POST",
            beforeSend: function () {
                if (commitMade) {
                    commitCall.abort();
                }
                commitMade = true;
            },
            success: function (data){
                // Hide the spinner
                $('#spinner').hide();
                commitMade = false;
                window.location.reload();
                $('#pageContent').show()
                $('#formContent').hide()
                displayWarningMessage(sitesCommittedSuccess)
            },

            error: function(data){
                console.log("error", data.responseText);
                // Hide the spinner if an error is returned
                $('#spinner').hide();
                commitMade = false;
                displayWarningMessage(data.responseText)
            }
        })
        }, function (){modalShown = false;})
}

