var siteOffset = 0;
var chainOffset = 0;
var acquirerOffset = 0;
var tidOffset = 0;
var othersOffset = 0;
var historyOffset = 0;


var After = ""
var Name = ""
var User = ""
var Before = ""
var Field = ""
var Offset = 0

function pageResults(pageAmount, tabIdent){
    if(tabIdent == "Site" && (pageAmount >= 0 || siteOffset > 0)){
        siteOffset = siteOffset + pageAmount; 
        filterRequest(tabIdent, siteOffset)
    } else if (tabIdent == "Chain" && (pageAmount >= 0 || chainOffset > 0)){
        chainOffset = chainOffset + pageAmount;
        filterRequest(tabIdent, chainOffset)
    } else if(tabIdent == "Acquirer" && (pageAmount >= 0 || acquirerOffset > 0)){
        acquirerOffset = acquirerOffset + pageAmount;
        filterRequest(tabIdent, acquirerOffset)
    } else if(tabIdent == "TID" && (pageAmount >= 0 || tidOffset > 0)){
        tidOffset = tidOffset + pageAmount;
        filterRequest(tabIdent, tidOffset)
    }else if(tabIdent == "Others" && (pageAmount >= 0 || othersOffset > 0)) {
        tidOffset = othersOffset + pageAmount;
        filterRequest(tabIdent, othersOffset)
    } else if(tabIdent == "Name" && (pageAmount >= 0 || historyOffset > 0)){
        historyOffset = historyOffset + pageAmount;
        filterHistoryRequest(tabIdent, historyOffset)
    } else{
        return;
    }
}

function getOffset (tabIdent) {
    switch (tabIdent) {
        case 'Site':
            return siteOffset
        case 'Chain':
            return chainOffset
        case 'Acquirer':
            return acquirerOffset
        case 'TID':
            return tidOffset
        case 'Others':
            return tidOffset
        case 'Name':
            return historyOffset
        default:
            return -1
    }
}

function resetOffset (tabIdent) {
    switch (tabIdent) {
        case 'Site':
            siteOffset = 0
            break
        case 'Chain':
            chainOffset = 0
            break
        case 'Acquirer':
            acquirerOffset = 0
            break
        case 'TID':
            tidOffset = 0
            break
        case 'Others':
            othersOffset = 0
            break
        case 'Name':
            historyOffset = 0
            break
        default:
            console.log("resetOffset() called with invalid tab identifier")
            break
    }
}

function filterRequest(tabIdent, offset){
    $.ajax({
        data: {
            Offset: offset,
            Type: tabIdent,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: "POST",
        url: "/filterChangeApproval",
        success: function(data){
            $("#changeApprovalContainer"+tabIdent).html(data);
            $("#shown-results-"+tabIdent).html("Showing Results "+offset+" to "+(offset+50));
            $("#shown-results-footer-"+tabIdent).html("Showing Results "+offset+" to "+(offset+50));
            if(offset > 0){
                $("#prev-results-"+tabIdent).removeClass("opacity");
                $("#prev-results-"+tabIdent).addClass("clickable");
                $("#prev-results-footer-"+tabIdent).removeClass("opacity");
                $("#prev-results-footer-"+tabIdent).addClass("clickable");
            } else {
                $("#prev-results-footer-"+tabIdent).removeClass("clickable");
                $("#prev-results-footer-"+tabIdent).addClass("opacity");
            }
        },
        error: function(data){
            displayWarningMessage(data.responseText)
        }
    })
}

function approveAll(profileType) {
    $.ajax({
        data: {
            Offset: getOffset(profileType),
            Type: profileType,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        url: "/approveAllChanges",
        success: function() {
            resetOffset(profileType)
            filterRequest(profileType, 0)
            resetBadges(profileType)
        },
        error: function(data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function approveChange(profile_data_id, profileType) {
    var data = {
        profileDataID: profile_data_id,
        Type: profileType
    }

    $.ajax({
        url: "/approveChange",
        data: data,
        success: function() {
            filterRequest(profileType, 0)
            updateBadge(profileType)
        },
        error: function(data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function discardAll(profileType) {
    $.ajax({
        data: {
            Offset: getOffset(profileType),
            Type: profileType,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        url: "/discardAllChanges",
        success: function() {
            filterRequest(profileType, 0)
            resetBadges(profileType)
        },
        error: function(data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function discardChange(profile_data_id, profileType) {
    var data = {
        profileDataID: profile_data_id,
        Type: profileType
    }
    $.ajax({
        url: "/discardChange",
        data: data,
        success: function() {
            filterRequest(profileType, 0)
            updateBadge(profileType)
        },
        error: function(data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function initComponents() {
    $("#after").datetimepicker({step: 5 })
    $("#before").datetimepicker({step: 5 })

    $("#clear-filters").click(function(){
        $('#before')[0].value = ""
        $('#after')[0].value = ""
        $('#name')[0].value = ""
        $('#user')[0].value = ""
        $('#field')[0].value = ""
        historyTabSelect("item", 0)
    })
}

function filterHistoryRequest(tabIdent, offset){
        After = $("#after")[0].value
        Name = $("#name")[0].value
        User = $("#user")[0].value
        Before = $("#before")[0].value
        Field = $("#field")[0].value

    $.ajax({
        data: {
            After: After,
            Name: Name,
            User: User,
            Before: Before,
            Field: Field,
            Offset: offset,
            Type: tabIdent,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: "POST",
        url: "/filterChangeApprovalHistory",
        success: function(data){
            $("#changeApprovalContainerHistory").html(data)
            $("#shown-results-"+tabIdent).html("Showing Results "+offset+" to "+(offset+50));
            $("#shown-results-footer"+tabIdent).html("Showing Results "+offset+" to "+(offset+50));
            if(offset > 0){
                $("#prev-results-"+tabIdent).removeClass("opacity");
                $("#prev-results-"+tabIdent).addClass("clickable");
                $("#prev-results-footer-"+tabIdent).removeClass("opacity");
                $("#prev-results-footer-"+tabIdent).addClass("clickable");
            } else {
                $("#prev-results-"+tabIdent).removeClass("clickable");
                $("#prev-results-"+tabIdent).addClass("opacity");
                $("#prev-results-footer-"+tabIdent).removeClass("clickable");
                $("#prev-results-footer-"+tabIdent).addClass("opacity");
            }

            initComponents()

            $("#after")[0].value = After
            $("#name")[0].value = Name
            $("#user")[0].value = User
            $("#before")[0].value = Before
            $("#field")[0].value = Field
        },
        error: function(data){
            displayWarningMessage(data.responseText)
        }
    })
}

function historyTabSelect(tabIdent, offset) {
    $.ajax({
        data: {
            Offset: offset,
            Type: tabIdent,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: "POST",
        url: "/filterChangeApprovalHistory",
        success: function(data){
            $("#changeApprovalContainerHistory").html(data)
            $("#shown-results-"+tabIdent).html("Showing Results "+0+" to "+50)
            if(offset > 0){
                $("#prev-results-"+tabIdent).removeClass("opacity")
                $("#prev-results-"+tabIdent).addClass("clickable")
            } else {
                $("#prev-results-"+tabIdent).removeClass("clickable")
                $("#prev-results-"+tabIdent).addClass("opacity")
            }
            initComponents()
        },
        error: function(data){
            displayWarningMessage(data.responseText)
        }
    })
}

function exportHistory() {
    let params = "After=" + After + "&" + "Name=" + Name + "&" + "User=" + User + "&" + "Before=" + Before + "&" + "Field=" + Field + "&" + "After=" + After + "&" + "Offset=" + Offset

    req.open("Get", "/exportChangeApprovalHistory?" + params, true);
    req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

    req.responseType = "blob";
    req.onload = function (event) {
        var blob = req.response;
        var fileName = req.getResponseHeader("fileName") //if you have the fileName header available
        var link=document.createElement('a');
        link.href=encodeURI(window.URL.createObjectURL(blob));
        link.download=fileName;
        link.click();
    };
    req.send(params);
};

function checkKey() {
    let key = event.which;
    if(key == 13) {
        filterHistoryRequest("Item", 0)
    }
}

function updateBadge(profileType) {
    if(profileType == "Site") {
        $("#siteBadge")[0].innerHTML = $("#siteBadge")[0].innerHTML - 1
        if($("#siteBadge")[0].innerHTML == 0) {
            $("#siteBadge")[0].setAttribute("hidden", "")
        }
    } else if(profileType == "Chain") {
        $("#chainBadge")[0].innerHTML = $("#chainBadge")[0].innerHTML -1
        if($("#chainBadge")[0].innerHTML == 0) {
            $("#chainBadge")[0].setAttribute("hidden", "")
        }
    } else if(profileType == "Acquirer") {
        $("#acquirerBadge")[0].innerHTML = $("#acquirerBadge")[0].innerHTML - 1
        if($("#acquirerBadge")[0].innerHTML == 0) {
            $("#acquirerBadge")[0].setAttribute("hidden", "")
        }
    } else if(profileType == "TID") {
        $("#tidBadge")[0].innerHTML = $("#tidBadge")[0].innerHTML - 1
        if($("#tidBadge")[0].innerHTML == 0) {
            $("#tidBadge")[0].setAttribute("hidden", "")
        }
    } else if(profileType == "Others") {
        $("#othersBadge")[0].innerHTML = $("#othersBadge")[0].innerHTML - 1
        if($("#othersBadge")[0].innerHTML == 0) {
            $("#othersBadge")[0].setAttribute("hidden", "")
        }
    }
}

function resetBadges(profileType) {
    if(profileType == "Site") {
        $("#siteBadge")[0].innerHTML = 0
        $("#siteBadge")[0].setAttribute("hidden", "")
    } else if(profileType == "Chain") {
        $("#chainBadge")[0].innerHTML = 0
        $("#chainBadge")[0].setAttribute("hidden", "")
    } else if(profileType == "Acquirer") {
        $("#acquirerBadge")[0].innerHTML = 0
        $("#acquirerBadge")[0].setAttribute("hidden", "")
    } else if(profileType == "TID") {
        $("#tidBadge")[0].innerHTML = 0
        $("#tidBadge")[0].setAttribute("hidden", "")
    } else if(profileType == "Others") {
        $("#othersBadge")[0].innerHTML = 0
        $("#othersBadge")[0].setAttribute("hidden", "")
    }
}

function downloadFlaggingFile(fileName){
    // filename will have `file : ` as well, So removing that
    fileName = fileName.substring(7)
    const data = {
        FileName: fileName,
        FileType: "TerminalFlagging",
        csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    }

    $.ajax({
        url: 'downloadFile',
        data: data,
        method: 'POST',
        success: (d, status, req) => {
            const link = document.createElement('a')
            link.href = encodeURI(window.URL.createObjectURL(new Blob([d], {type: 'text/plain'})));
            // filename will have date as well, So removing the date from the file name (yyyymmddhhmmss_)
            link.download = fileName.substring(15);
            link.click()
        },
        error: (d) => {
            displayWarningMessage(d.responseText)
        }
    })
}

