function showConfirmDialog(title, msg){
    var confirm = $('#confirmModal')
    var p = jQuery.Deferred()
    $('#confirm-title').html(title)
    $('#confirm-message').html(msg)
    $('#confirm-yes').click(function(){
        confirm.modal('hide')
        p.resolve(true)
        return true
    })
    $('#confirm-no').click(function(){
        confirm.modal('hide')
        p.resolve(false)
        return false
    })

    confirm.modal('show')

    return p.promise()
}

function confirmDialog(title, msg, onSuccess){
    var p = showConfirmDialog(title,msg)
    p.done( function(confirmed){
        if(confirmed){
            onSuccess()
        }
    })

}