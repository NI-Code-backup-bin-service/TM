var profileType = "site"
var hiddenDataGroups = false

function toggleDataGroups(){
    if(hiddenDataGroups){
        $("#data-groups").removeClass("hidden")
        $("#data-groups").addClass("visible")
        $("#dgToggle").html("hide")
    } else {
        $("#data-groups").removeClass("visible")
        $("#data-groups").addClass("hidden")
        $("#dgToggle").html("show")
    }

    hiddenDataGroups = !hiddenDataGroups;
}


function getDataGroups(){
    var data = {
        AcquirerId: $("#acquirer_dropdown option:selected").val(),
        ChainId: $("#chain_dropdown option:selected").val(),
        UseGlobal: "1",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }

    if(profileType == "acquirer"){
        data.ChainId = "-1"
        data.AcquirerId = "-1"      
    }

    if(profileType == "chain") {
        data.ChainId = "-1"
    }

    $.ajax({ 
        url: "getDataGroups",
        data: data,
        method: "POST",
        success: function (d) {
            $('#data-groups').html(d);
        },
        error: function(data){
            
        }
    });
}
function typeSelect(type){
        if (type.value == "site")
        {
            $("#chain_select").show()
            $("#acquirer_select").hide()
            $("#addSiteOptions").show()
            $("#addProfileOptions").hide()
        }
        else if (type.value == "chain")
        {
            $("#chain_select").hide()
            $("#acquirer_select").show()
            $("#addSiteOptions").hide()
            $("#addProfileOptions").show()
        }
        else
        {
            $("#chain_select").hide()
            $("#acquirer_select").hide()
            $("#addSiteOptions").hide()
            $("#addProfileOptions").show()
        }  

    profileType = type.value;

    getDataGroups()

}

$(document).ready(function(){
    $('#save-new-site').submit( function(e) {
        const saveButton = $("#site-save")[0];
        // We disable the button when the form is submitted to stop users from pressing it multiple times which would
        // cause multiple site creation requests to be sent. We don't need to enable it again in the success of the ajax
        // request as we direct to  the search page.
        saveButton.disabled = true;
        e.preventDefault();

        var acquirerId =  $("#acquirer_dropdown option:selected").val();
        var chainId =  $("#chain_dropdown option:selected").val();

        var formData = $(this).serializeArray()
        formData.push({name: "type", value: profileType})
        formData.push({name: "name", value: $("#name").val()})
        formData.push({name: "acquirer", value: acquirerId })
        formData.push({name: "chain", value: chainId})
        $.ajax({
            data: formData,
            typeField: profileType,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            success: function(data){
                window.location.href = "/search"
            },
            error: function(data){
                saveButton.disabled = false;
                //the data is a json array of messages
                if(data.getResponseHeader("content-type") === "application/json") {
                    let validationMessages = JSON.parse(data.responseText).join('<br/>');
                    displayWarningMessage(validationMessages)
                } else {
                    displayWarningMessage(data.responseText)
                }
            }

        })
    })

    $('#get-fields').submit( function(e) {

        e.preventDefault();
        var formData = $(this).serializeArray()
        var dgData = serializeAllDataGroups();
        $.merge(formData, dgData)
        formData.push({name: "type", value: profileType})
        $.ajax({
            data: formData,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            success: function(d){
                $('#addSiteFields').html(d);
                bindMultiSelect();
                bindOverrideField();
                hideWarning();
            },
            error: function(data){
                if(data.getResponseHeader("content-type") === "application/json") {
                    let validationMessages = JSON.parse(data.responseText).join('<br/>');
                    displayWarningMessage(validationMessages)
                } else {
                    displayWarningMessage(data.responseText)
                }
            }

        })
    })

    $("#addSiteOptions").show()
    $("#acquirer_select").hide()
    $("#addProfileOptions").hide()
})

// This function is a workaround to allow us to get all the data groups that are enabled, this includes the
// disabled but checked ones.
function serializeAllDataGroups() {
    const dataGroupForm = $("#save-new-site");
    let disabledGroups = dataGroupForm.find(':input:disabled').removeAttr('disabled');
    let serializedForm = $("[id^=dg]").serializeArray();
    disabledGroups.attr('disabled','disabled');
    return serializedForm;
}