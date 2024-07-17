var updateId = 0;

function searchUpdatesSN(tidId, siteId, serialNo) {
   $(".xdsoft_datetimepicker").detach()
    let data = {
        TID: tidId,
        Site: siteId,
        SerialNo: serialNo,
        Serial: "SN",
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    $.ajax({
        url: "updatesSN",
        data: data,
        method: "POST",
        success: function (d) {
            let clean = sanitizeHTML(d)
            $( "#snUpdatesModalBody" ).html(clean);
            bindDatePickers();
            $("#snUpdatesModal").modal('show');
            updateId = parseInt($("#update_table_sn tr:last").find('td:first').attr('version-select-sn'));
            if ( isNaN(updateId) ) { updateId = 0; }
            document.getElementById("applyUpdatesSN").addEventListener("click", function(event){
                event.preventDefault()
            });
            let snEditInput = "#SnEditInput"+tidId
            let snCell = "#SnCell"+tidId
            $(snEditInput).attr('value', serialNo)
            $(snCell).attr('value', serialNo)
            $("#applyUpdatesSN").attr("disabled", "disabled")
            $("#update_table_sn").find("input,select").attr("disabled", "disabled");

            updateId = $('#update_table_sn').find('tr').length
            let addTid = new FormData($("#snUpdates_" + tidId)[0]);
            addTid.append("TID", tidId);
            addTid.append("siteId", siteId);
            addTid.append("Serial", "SN");
            addTid.append("UpdateID", ++updateId);
            addTid.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value);
            $.ajax({
                type: "POST",
                data: addTid,
                url: "/addTIDUpdate",
                processData: false,
                contentType: false,
                success: function(tidData) {
                    $("#update_table_sn").append(sanitizeTableHTML(tidData))
                    bindDatePickers()
                    let sid = "#version_select_sn_"+updateId
                    $(sid+' option:selected').not(':disabled').removeAttr('selected');
                    document.getElementsByName(sid).selectedIndex = "0";

                },
                error: function(data) {
                    displayWarningMessage(data.responseText,"updateDetailsSNPartial")
                }
            })
        },        
    });
}


function SetDirty() {
    $("#applyUpdatesSN").removeAttr("disabled")
}

function ApplyTIDSNUpdates(tidId) {
    $(".xdsoft_datetimepicker").remove();
    let updatedSn = $("#SnEditInput"+tidId).val();
    let oldSn = $("#SnCell"+tidId).val();
    if (updatedSn == oldSn) {
        $('#snUpdatesModal').modal('hide');
    }
    let apkIDs = []
    let apkID = $("#apkIDs").val();
    if (apkID != null) {
        apkIDs = apkID;
    }
    let updateID = 0;
    let tidUpdateID = $("#tidUpdateID").val();
    if  (tidUpdateID != null){
        updateID = tidUpdateID
    }
    let dataObj = new FormData($("#snUpdates_" + tidId)[0]);
    dataObj.append("TID", tidId);
    dataObj.append("OLD_SN", oldSn);
    dataObj.append("NEW_SN", updatedSn);
    dataObj.append("SITE_PROFILE_ID", $("#profileID").val());
    dataObj.append("thirdParty_"+updateID, apkIDs);
    dataObj.append('csrfmiddlewaretoken', $("input[name=csrfmiddlewaretoken]")[0].value);
    $.ajax({
        type: "POST",
        url: "/updateSerialNumber",
        data: dataObj,
        processData: false,
        contentType: false,
        success: function (data) {
            $("#update_table_sn").append(sanitizeTableHTML(data))
            $("#applyUpdatesSN").attr("disabled", "disabled")
            location.reload();
        },
        error: function (data) {
            displayWarningMessage(data.responseText,"updateDetailsSNPartial")
        }
    })
}
