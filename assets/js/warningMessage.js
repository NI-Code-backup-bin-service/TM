

function displayWarningMessage(message, ident){
    if (ident === undefined) {
        ident = "";
    }

    $("#warning-message_"+ident)[0].innerHTML = message;

    $("#warning_"+ident).show();
}

function hideWarning(ident){
    if (ident === undefined) {
        ident = "";
    }
    $("#warning_"+ident).hide()
}