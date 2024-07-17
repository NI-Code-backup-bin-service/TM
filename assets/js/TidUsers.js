function saveTidUsers(usersList, deletedSet, modfiedSet) {
    const modifiedUsers = modfiedSet;
    let newUsers = [];
    let users = [];
    newUsers = usersList.filter(function (item) {
        // Check if UserId is -1 and Username is not empty
        if (item.UserId < 0  && item.Username.trim() !== "") {
            // Remove UserId key if it is -1
            if (item.UserId < 0) {
                delete item.UserId;
            }
            return true; // Include this user in the filtered array
        }
        //if there is no userId e.g. undefined 
        if (item.Username.trim() !== ""){
            return true ; // Exclude this user from the filtered array
        }
        return false; // Exclude this user from the filtered array
    });

    users = usersList.filter(function (item) {
        // Check if UserId is in modifiedUsers set, and Username is not empty
        if ( modifiedUsers.has(item.UserId) && item.Username.trim() !== "") {
            return true; // Include this user in the filtered array
        }
        return false; // Exclude this user from the filtered array
    });
    let data = new FormData()
    
    const deletedUsers = Array.from(deletedSet);
    data.append("Users",  JSON.stringify(users));
    data.append("NewUsers",  JSON.stringify(newUsers));
    data.append("DeletedUsers",  JSON.stringify(deletedUsers));
    data.append("SiteId",  $("#siteID")[0].value);
    data.append("TID",  $("#UserTid")[0].value);
    data.append("profileID",  $("#profileID")[0].value);
    data.append("csrfmiddlewaretoken",  $("input[name=csrfmiddlewaretoken]")[0].value);

    $.ajax({
        url: "/SaveTidUsers",
        data: data,
        method: "POST",
         processData: false,
        contentType: false,
        success: function (d) {
            $("#tidUsersModal").modal('hide');
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised","tidUserOverride")
            } else {
                displayWarningMessage(data.responseText,"tidUserOverride")
            }
        }
    });
}

function clearUsers(tid){
    confirmDialog("Clear Overrides?",'Are you sure you want to remove the user overrides for tid: '+tid+' ?', function(){
        const data = {
            TID:  $("#UserTid")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        };
        $.ajax({
            url: "/ClearTidUsers",
            data: data,
            method: "POST",
            success: function (d) {
                $("#tidUsersModal").modal('hide')
                $("#user-override-"+$("#UserTid")[0].value).removeClass("activated-btn")
            },
            error: function (data) {
                if(data.status == 401 || data.status == 403 ){
                    displayWarningMessage("User not authorised","tidUserOverride")
                } else {
                    displayWarningMessage(data.responseText,"tidUserOverride")
                }
            }
        });
    });
}

function fetchTidUsers(){
    const tid = $("#UserTid")[0].value

    $.ajax({
        data: {
            siteId: $("#siteID")[0].value,
            TID:  tid,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: "POST",
        url: "/GetTidUsers",

        success: function (d) {
            $('#tidUsersModalTitle').html("TID users (" + tid + ")");
            if(d.Users === null){
                d.Users = [];
            }
            d.Users.push({UserId: addedRowId--, Username: '', PIN: generateNewRandomPin(d.Users), Modules: ['sale','gratuitySale','X-Read','Z-Read']});
            const tidUserTable = ce(UserTable, {
                key:"tid-user-table",
                Modules: d.Modules,
                Users: d.Users,
                FriendlyModules: d.FriendlyModules,
                hasSavePermission: d.HasSavePermission,
                hasPasswordPermission: d.HasPasswordPermission,
                UserSave: saveTidUsers,
                PerPage: 7,
                DisplayWarningMessage: showModalWarning
            }, null);

            ReactDOM.unmountComponentAtNode(document.getElementById("tidUserTable"));
            ReactDOM.render(tidUserTable, document.getElementById("tidUserTable"));

            const userOverride = $('#user-override-' + tid)
            if (d.Overriden) {
                $('#clear-tid-overrides').removeClass('hidden')
                userOverride.addClass('activated-btn')
            } else {
                userOverride.removeClass('activated-btn')
            }
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised","tidUserOverride")
            } else {
                displayWarningMessage(data.responseText,"tidUserOverride")
            }
        }
    })
}

function showModalWarning(message){
    displayWarningMessage(message, "tidUserOverride")
}

$(document).ready(function () {
    fetchTidUsers();
});