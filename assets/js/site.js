let req = new XMLHttpRequest();
let socket

$(document).ready(function() {
    $("#searchForm").submit(function (event) {
        event.preventDefault();
        var encodedSearchTerm = encodeURIComponent($("#searchTerm").val());
        var action = $(this).attr('action')
        if ( $( "#searchPartial" ).length ) {
            $.post(action, $(this).serialize()).done(function (data) {
                $('#searchPartial').html(sanitizeHTML(data));
                bindEditProfile();
            });;
        }else{
            $(location).attr('href', action + "?searchTerm=" + encodedSearchTerm)
        }
    });

    bindEditProfile();
    bindDeleteSite();
    bindDeleteChain();
    bindDeleteAcquirer();

    socket = new WebSocket('wss://' + window.location.hostname + ':5006/echo')
    bindSocket()

    // Select all tabs
    $('.nav-tabs a').click(function(){
        $(this).tab('show');
    })

    $('.nav-pills a').click(function(){
        $(this).tab('show');
    })

    $("#site").tab('show');
    $("#chain").tab('show');

    if( ('undefined' !== typeof navOpen) && navOpen) {
        openNav()
    }
    else {
        closeNav()
    }
});

function bindSocket(){
    //Redirects the user to logon when session expired
    socket.onmessage = function (e) {
        window.location.replace("/signon")
    }
}


function bindEditProfile(){
    $("[data-button='edit-site']").click(function(){
        $(location).attr('href', '/profileMaintenance?profileId=' + $(this).attr("data-site"))+'&type=site'
    })

    $("[data-button='edit-chain']").click(function(){
        $(location).attr('href', '/profileMaintenance?profileId=' + $(this).attr("data-chain"))+'&type=chain'
    })

    $("[data-button='edit-acquirer']").click(function(){
        $(location).attr('href', '/profileMaintenance?profileId=' + $(this).attr("data-acquirer"))+'&type=acquirer'
    })
}

function bindDeleteSite(){
    $("[data-button='delete-site']").click(function(){
        var site = $(this).attr("data-site")
        confirmDialog("Delete Site?",'Are you sure you want to delete Site?', function(){
            var data = {
                siteId: site,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            }
        
            $.ajax({ 
                url: "deleteSite", 
                data: data,
                method: "POST",
                success: function (d) {
                    location.reload()
                },
                error: function(data){
                    displayWarningMessage(data.responseText)
                }
            });;
        })
    })   
}

function bindDeleteChain(){
    $("[data-button='delete-chain']").click(function(){
        var chain = $(this).attr("data-chain")
        confirmDialog("Delete Chain?",'Are you sure you want to delete Chain?', function(){
            var data = {
                chainId: chain,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            }

            $.ajax({
                url: "deleteChain",
                data: data,
                method: "POST",
                success: function (d) {
                    location.reload()
                },
                error: function(data){
                    displayWarningMessage(data.responseText)
                }
            });;
        })
    })
}

function bindDeleteAcquirer(){
    $("[data-button='delete-acquirer']").click(function(){
        var acquirer = $(this).attr("data-acquirer")
        confirmDialog("Delete Acquirer?",'Are you sure you want to delete Acquirer?', function(){
            var data = {
                acquirerId: acquirer,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            }

            $.ajax({
                url: "deleteAcquirer",
                data: data,
                method: "POST",
                success: function (d) {
                    location.reload()
                },
                error: function(data){
                    displayWarningMessage(data.responseText)
                }
            });
        })
    })
}

function openNav() {
    document.getElementById("sideNav").style.width = "250px";
    document.getElementById("main").style.marginLeft = "250px";
}

function closeNav() {
        document.getElementById("sideNav").style.width = "0";
        document.getElementById("main").style.marginLeft= "0";
}

function exportTooltips() {
    let ttModal = $('#tooltipExportModal')
    ttModal.modal({ backdrop: 'static', keyboard: false})
    ttModal.modal('show')

    $.ajax({
        url: '/exportTooltips',
        success: (d, status, req) => {
            const fileName = req.getResponseHeader('fileName') //if you have the fileName header available
            const link = document.createElement('a')
            link.href = encodeURI(window.URL.createObjectURL(new Blob([d], {type: 'text/csv'})));
            link.download = fileName
            link.click()
            $('#tooltipExportModal').modal('hide')
        },
        error: (d) => {
            $('#tooltipExportModal').modal('hide')
            setTimeout(function () {
                alert(d.responseText)
            }, 100);
        }
    })
}

function backupDatabase() {
    $("#DbBackupModal").modal( {backdrop: 'static', keyboard: false} )
    $("#DbBackupModal").modal('show')

    let params = ""
    req.open("Get", "/backupDatabase?" + params, true);
    req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    req.responseType = "blob";
    req.onload = function (event) {
        var blob = req.response;
        var fileName = req.getResponseHeader("fileName") //if you have the fileName header available
        var link=document.createElement('a');
        link.href=encodeURI(window.URL.createObjectURL(blob));
        link.download=fileName;
        link.click();
        $("#DbBackupModal").modal('hide')
    };
    req.send(params);
}

