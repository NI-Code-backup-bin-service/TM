var superPins = [];
var updateUsersResponse = null;

function exportCsv() {
    const data = {
        SiteId: $("#siteID")[0].value,
        ProfileId: $("#profileID")[0].value,
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };
    $.ajax({
        url: "/ExportUserCsv",
        data: data,
        method: "POST",
        success: function (d) {
            const csvFile = new Blob([d], {type: "text/csv"});
            let link=document.createElement('a');
            link.href=encodeURI(window.URL.createObjectURL(csvFile));
            link.download="Users.csv";
            link.click();
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised")
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    });
}

function saveSiteUsers(usersList, deletedSet, modfiedSet) {
    const modifiedUsers = modfiedSet;
    let flagStatus = $("#hidden_userFlagStatusOption").val();
    let users = [];
    let newUsers = [];
    if (flagStatus != "") {
        newUsers = usersList.filter(function (item) {
            // Check if UserId is -1 and Username is not empty
            if (item.UserId < 0  && item.Username.trim() !== "") {
                // Remove UserId key if it is -1
                if (item.UserId < 0) {
                    delete item.UserId;
                }
                return true; // Include this user in the filtered array
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
    } else {
        displayWarningMessage("You must select a flagging option")
        return
    }

    const deletedUsers = Array.from(deletedSet);
    let dataObj = new FormData()
    if (flagStatus=="file"){
        var uploadFile = $('#flagged-tids-file')[0].files[0];
        dataObj.append("flagged-tids-file",  uploadFile);
    }
    dataObj.append("Users",  JSON.stringify(users));
    dataObj.append("NewUsers",  JSON.stringify(newUsers));
    dataObj.append("DeletedUsers",  JSON.stringify(deletedUsers));
    dataObj.append("SiteId",  $("#siteID")[0].value);
    dataObj.append("profileID",  $("#profileID")[0].value);
    dataObj.append("flagStatus",  flagStatus);
    dataObj.append("csrfmiddlewaretoken",  $("input[name=csrfmiddlewaretoken]")[0].value);

    $("#flagging-checkbox").find("input:checkbox:checked").not("[disabled]").each(function() {
        dataObj.append($(this).attr('id'),$(this).val());
    });

    $.ajax({
        url: "/SaveSiteUsers",
        data: dataObj,
        method: "POST",
        processData: false,
        contentType: false,
        success: function (d) {
            displayWarningMessage("Saved")
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised")
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    });
}

function fetchUsers(){
    $.ajax({
        data: {
            siteId: $("#siteID")[0].value,
            profileId: $("#profileID")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: "POST",
        url: "GetSiteUsers",

        success: function (d) {
            superPins = d.SuperPins;
            if(d.Users === null){
                d.Users = [];
            }
            d.Users.push({UserId: addedRowId--, Username: '', PIN: generateNewRandomPin(d.Users), Modules: ['sale','gratuitySale','X-Read','Z-Read']});
            const siteUserTable = ce(UserTable, {
                key:"site-user-table",
                Modules: d.Modules,
                Users: d.Users,
                FriendlyModules: d.FriendlyModules,
                PerPage: 10,
                hasSavePermission: d.HasSavePermission,
                hasPasswordPermission: d.HasPasswordPermission,
                UserSave: saveSiteUsers,
                DisplayWarningMessage: displayWarningMessage
            }, null);
            ReactDOM.unmountComponentAtNode(document.getElementById("siteUserTable"));
            ReactDOM.render(siteUserTable, document.getElementById("siteUserTable"));
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised")
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    })
}

$(document).ready(function () {
    $('[data-bs-toggle="tooltip"]').tooltip();
    fetchUsers();
});

$(document).ready(function(){
    $('[data-bs-toggle="tooltip"]').tooltip();
    $('#upload-csv').submit( function(e) {
        e.preventDefault();
        var dataObj = new FormData(document.querySelector("#upload-csv"))
        dataObj.append("SiteId",  $("#siteID")[0].value)

        $.ajax({
            data: dataObj,
            type: $(this).attr('method'),
            url: $(this).attr("action"),
            processData: false,
            contentType: false,

            success: function(data){
                ShowUserUploadReport(data);
                fetchUsers();
            },
            error: function(data){
                displayWarningMessage(data.responseText)
            }

        })
    })
});

function ExportUserUploadReportAsCSV() {
    if(updateUsersResponse == null) {
        return false;
    }
    const data = {
        ProfileId: $("#profileID")[0].value,
        UploadResult: JSON.stringify(updateUsersResponse),
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    };

    $.ajax({
        url: "/ExportUploadUserCSVResult",
        data: data,
        method: "POST",
        success: function (d) {
            const csvFile = new Blob([d], {type: "text/csv"});
            let link=document.createElement('a');
            link.href=encodeURI(window.URL.createObjectURL(csvFile));
            link.download="UserUploadResult.csv";
            link.click();
        },
        error: function (data) {
            if(data.status == 401 || data.status == 403 ){
                displayWarningMessage("User not authorised")
            } else {
                displayWarningMessage(data.responseText)
            }
        }
    });
}

function ShowUserUploadReport(saveResult) {
    updateUsersResponse = saveResult;

    var numberOfAddedUsers = saveResult.filter(function (userResult) {
        return (userResult.Result.Success && userResult.Action.includes("Add"));
    }).length;

    var numberOfChangedUsers = saveResult.filter(function (userResult) {
        return (userResult.Result.Success && userResult.Action.includes("Update"));
    }).length;

    var numberOfDeletedUsers = saveResult.filter(function (userResult) {
        return (userResult.Result.Success && userResult.Action.includes("Delete"));
    }).length;

    var numberOfFailures = saveResult.filter(function (userResult) {
        return !userResult.Result.Success;
    }).length;

    var resultCountLabels = ""
    if (numberOfAddedUsers > 0) {
        resultCountLabels = resultCountLabels.concat(`<div class="alert alert-info" role="alert">Number of users added: ${numberOfAddedUsers}</div>`)
    }
    if (numberOfChangedUsers > 0) {
        resultCountLabels = resultCountLabels.concat(`<div class="alert alert-info" role="alert">Number of users updated: ${numberOfChangedUsers}</div>`)
    }
    if (numberOfDeletedUsers > 0) {
        resultCountLabels = resultCountLabels.concat(`<div class="alert alert-info" role="alert">Number of users deleted: ${numberOfDeletedUsers}</div>`)
    }
    if (numberOfFailures > 0) {
        resultCountLabels = resultCountLabels.concat(`<div class="alert alert-danger" role="alert">Number of failures: ${numberOfFailures}</div>`)
    }
    $('#update-summary').empty();
    $('#update-summary').append(resultCountLabels);
    $('#uploadResultModal').modal('show');
}