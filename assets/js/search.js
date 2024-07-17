var ActiveTab = '#sites'
var ProcessingExport
const merchantIDIndex = 5;
$(document).ready(function () {
  search()
  ProcessingExport = $('#userData').attr('data-user')
  if (ProcessingExport !== '') {
    let exportModal = $('#exportModal')
    exportModal.modal({ backdrop: 'static', keyboard: false })
    exportModal.modal('show')
  }
})

function bindToggle () {
  $('[data-button="toggle"]').click(function () {
    $(this).parents().next('.hide').toggle()
  })

  $('[data-button="toggle"]').parents().next('.hide').toggle()
}

function checkKey () {
  const key = event.which
  if (key === 13) {
    search()
  }
}

function setDataTablesError () {
  $.fn.dataTable.ext.errMode = function ( settings, helpPage, message ) {
    location.reload()
  }
}

function search () {
  const token = $('input[name=csrfmiddlewaretoken]')[0].value
  $.ajax({
    url: 'search',
    method: 'GET',
    success: function (d) {
      setDataTablesError()
      bindDeleteChain()
      bindDeleteAcquirer()
      bindToggle()
      buildSitesTable(token)
      buildTIDsTable(token)
      buildChainsTable(token)
      buildAcquirersTable(token)
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}


function addSearchKeyBind(searchInputId, searchButtonId) {

  if ($.type(searchInputId) !== 'string' || $.type(searchButtonId) !== 'string') {
    console.log('Invalid arguments for function addSearchKeyBind');
    return;
  }

  $('#' + searchInputId).keyup(function(event){
    if(event.key !== "Enter") return;
    event.preventDefault();
    $('#' + searchButtonId).click();
  });
}

function buildSitesTable (token) {
  let table = $('#sitesTable').DataTable({
    ajax: {
      'url': location.href,
      'type': 'POST',
      'data': (d) => {
        return $.extend({}, d, {
          'requestType': 'site',
          csrfmiddlewaretoken: token
        })
      }
    },
    deferRender: true,
    'processing': true,
    "deferLoading": 0,
    'serverSide': true,
    'pageLength': 15,
    'lengthMenu': [
      [10, 15, 25, 50, 100, 500, -1],
      [10, 15, 25, 50, 100, 500, 'All']
    ],
    'columns': [
      { data: 'MerchantId' },
      { data: 'SiteName' },
      { data: 'ChainName' },
      {
        data: null,
        'orderable': false,
        'render': (cellData, cellType, rowData, meta) => {
          const mid = rowData['MerchantId']
          return '<div class="btn-group btn-group-sm" role="group"> <button id="edit-site-' + mid + '" type="button" class="btn btn-secondary btn-sm button-margin" ' +
            'data-button="edit-site">Edit' +
            '</button> ' +
            '<button id="delete-site-' + mid + '" type="button" class="btn btn-outline-danger btn-sm button-margin" ' +
            'data-button="delete-site"' + (permSiteDelete ? '' : ' disabled') + '>Delete' +
            '</button> </div>'
        }
      }
    ],
    'autoWidth': false,
    'createdRow': (row, data, index) => {
      $(row).attr('id', 'site-row-' + data['SiteID'])
    },
    'columnDefs': [{
      'targets': '_all',
      'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
        $(td).attr('id', 'site-cell-' + colIndex + "-" + rowData['SiteID'])
      }
    }],
    "drawCallback": ( settings ) => {
      setPageIds()
    },
    initComplete: function() {
      // Use of unbind() means that the search input will no longer auto-filter. This is required because auto search
      // reloads the /searchHandler which then does a load of DB calls and negatively impacts performance
      let input = $('#sitesTable_filter input').attr('id', 'searchTerm').unbind();
      let self = this.api();

      // User must click this button to search now
      let $searchButton = $('<button>')
          .text('Search')
          .attr('id', 'siteSearchButton')
          .attr('class', 'btn btn-primary button-margin')
          .click(function() {
            self.search(input.val()).draw();
          });

      $('#sitesTable_filter').append($searchButton);

      // Add Event trigger to execute search when enter key is pressed
      addSearchKeyBind('searchTerm', 'siteSearchButton');
    }
  });

  $('#sitesTable tbody').on('click', 'button', function () {
    let row = table.row($(this).parents('tr'))
    let data = row.data()
    const profID = data["SiteProfileID"]

    if ($(this).attr('data-button') === 'edit-site') {
      location.href = '/profileMaintenance?profileId=' + profID + '&type=site'
      return
    }

    confirmDialog('Delete Site?', 'Are you sure you want to delete Site?', () => {
      $.ajax({
        url: 'deleteSite',
        data: {
          siteId: profID,
          csrfmiddlewaretoken: token
        },
        method: 'POST',
        success: function (d) {
          $('#site-row-' + data["SiteID"]).remove()
          hideWarning()
        },
        error: function (data) {
          displayWarningMessage(data.responseText)
        }
      })
    })
  })
}

function buildTIDsTable (token) {
  let table = $('#tidsTable').DataTable({
    ajax: {
      'url': location.href,
      'type': 'POST',
      'data': (d) => {
        return $.extend({}, d, {
          'requestType': 'tid',
          csrfmiddlewaretoken: token
        })
      }
    },
    deferRender: true,
    'processing': true,
    "deferLoading": 0,
    'serverSide': true,
    'pageLength': 15,
    'lengthMenu': [
      [10, 15, 25, 50, 100, 500, -1],
      [10, 15, 25, 50, 100, 500, 'All']
    ],
    order: [[1, 'asc']],
    'columns': [
      { data: 'TID',
        'render': (cellData, cellType, rowData) => {
            return rowData['TID'].toString().padStart(8, "0")
        }
      },
      { data: 'Serial' },
      { data: 'EnrolmentPIN' },
      { data: 'ResetPIN' },
      { data: 'ActivationTime' },
      { data: 'MerchantID' },
      { data: 'SiteName' },
      {
        data: null,
        'orderable': false,
        'render': (cellData, cellType, rowData, meta) => {
          //Pad the TID, we do this so that TIDs that start with 0 display correctly
          const tid = rowData['TID'].toString().padStart(8, "0")
          return '<button id="generate-enrolment-pin-' + tid + '" class="btn btn-secondary btn-sm button-margin" data-button="generate-enrolment" ' +
            (permSiteWrite ? '' : 'disabled') + '>Generate Enrolment PIN' +
            '</button>\n' +
            '<button id="generate-reset-pin-' + tid + '" class="btn btn-secondary btn-sm button-margin" data-button="generate-reset" ' +
            (permSiteWrite ? '' : 'disabled') + '>Generate Reset PIN' +
            '</button>\n' +
            '<button id="delete-' + tid + '" class="btn btn-secondary btn-sm button-margin" data-button="delete" ' +
            (permSiteWrite ? '' : 'disabled') + '>Delete' +
            '</button>\n' +
            '<button id="manage-' + tid + '" class="btn btn-secondary btn-sm button-margin" data-button="updates" ' +
            (permSiteWrite ? '' : 'disabled') + '>Manage Updates' +
            '</button>\n' +
            '<button id="detail-' + tid + '" class="btn btn-secondary btn-sm button-margin" data-button="details"' +
            '>Details' +
            '</button>'
        }
      }
    ],
    'autoWidth': false,
    'createdRow': (row, data, index) => {
      $(row).attr('id', 'tid-row-' + data['TID'].toString().padStart(8, "0"))
    },
    'columnDefs': [{
      'targets': '_all',
      'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
        if (colIndex == merchantIDIndex){
            $(td).attr('id', 'tid-cell-' + colIndex + "-" + rowData['TID'].toString().padStart(8, "0"))
                 .attr("profile-id", rowData['TIDProfileID'].toString())
                .attr("class","mid-highlight")
        }else{
            $(td).attr('id', 'tid-cell-' + colIndex + "-" + rowData['TID'].toString().padStart(8, "0"))
        }
      }
    }],
    "drawCallback": ( settings ) => {
      setPageIds()
    },
    initComplete: function() {
      let input = $('#tidsTable_filter input').attr('id', 'tids-searchTerm').unbind();
      let self = this.api();
      // User must click this button to search now
      let $searchButton = $('<button>')
          .text('Search')
          .attr('id', 'tidSearchButton')
          .attr('class', 'btn btn-primary button-margin')
          .click(function() {
            self.search(input.val()).draw();
          });
      $('#tidsTable_filter').append($searchButton);
      // Add Event trigger to execute search when enter key is pressed
      addSearchKeyBind('tids-searchTerm', 'tidSearchButton');
    }
  })

  $("#tidsTable tbody").on('click', 'td', function () {
      var profileID = $(this).attr('profile-id')
      typeof profileID !== 'undefined' ? location.href = '/profileMaintenance?profileId='+profileID+'&type=site': "";
  });

  $('#tidsTable tbody').on('click', 'button', function () {
    let row = table.row($(this).parents('tr'))
    let data = row.data()
    const tid = data['TID'].toString().padStart(8, "0")
    const siteID = data['SiteId']
    const serialNo = data['SerialNo']

    switch ($(this).attr('data-button')) {
      case 'generate-enrolment':
        GenerateOTP(tid, otpIntentEnum.Enrolment, (d) => {
          table.draw()
        })
        break
        case 'generate-reset':
          GenerateOTP(tid, otpIntentEnum.Reset, (d) => {
            table.draw()
          })
          break
      case 'delete':
        searchDeleteTID(tid, siteID, () => {
          alert("Deletion of TID " + tid + " added to change approval")
          table.draw()
        })
        break
      case 'updates':
        searchUpdatesTID(tid, siteID)
        break
      case 'updates_sn':
        searchUpdatesSN(tid, siteID, serialNo)
        break
      case 'details':
        $.ajax({
          url: "getTidDetails",
          data: {
            TID: tid,
            csrfmiddlewaretoken: token
          },
          method: "POST",
          success: (d) => {
            let clean = sanitizeHTML(d);
            $("#tidModalBody").html(clean)
            $("#tidModal").modal('show')
          },
          error: (d) => {
            displayWarningMessage(d.responseText)
          }
        })
        break;
    }
  })
}

function createDuplicateChain(data){
  let val = $("#newChainInput-"+data.profileId).val();
  var token =  $("input[name=csrfmiddlewaretoken]")[0].value;

  $.ajax({
    url: "addNewDuplicatedChain",
    data: {
      chainProfileId:data['profileId'],
      acquirerName:data['acquirerName'],
      newChainName:val,
      csrfmiddlewaretoken: token,
    },
    method: "POST",
    success: function (d) {
      $('#chainSearchButton').click();
      hideWarning("")
    },
    error: (d) => {
      displayWarningMessage(d.responseText);
    }
  });
}

function format(profileId,acquirerName) {
  data = {
    "profileId":profileId,
    "acquirerName":acquirerName
  }
  return (
      `<table style="padding-left:50px;">
      <tr>
      <td>Chain Name:</td>
      <td>
      <input id="newChainInput-`+profileId+`" name="newChaninInput" type="text" class="items col-md-auto">
      <input id="newChainacquireName-`+acquirerName+`" name="newChainacquireName" type="text" class="items col-md-auto hidden">
      </td>
      <td>
      <button id="save" type="button" class="items btn btn-primary btn-sm float-end" onclick='createDuplicateChain(`+JSON.stringify(data)+`)'>Apply Duplicate Chain</button>
      </td>
      </tr>
      </table>`
  );
}

function buildChainsTable (token) {
  let table = $('#chainsTable').DataTable({
    ajax: {
      'url': location.href,
      'type': 'POST',
      'data': (d) => {
        return $.extend({}, d, {
          'requestType': 'chain',
          csrfmiddlewaretoken: token
        })
      }
    },
    deferRender: true,
    'processing': true,
    "deferLoading": 0,
    'serverSide': true,
    'pageLength': 15,
    'lengthMenu': [
      [10, 15, 25, 50, 100, 500, -1],
      [10, 15, 25, 50, 100, 500, 'All']
    ],
    'columns': [
      { data: 'ChainProfileID' },
      { data: 'ChainName' },
      { data: 'ChainTIDCount' },
      { data: 'AcquirerName' },
      {
        data: null,
        'orderable': false,
        'render': (cellData, cellType, rowData, meta) => {
          const id = rowData['ChainProfileID']
          return '<button id="duplicate-chain-' + id + '" data-href="#content2" class="btn btn-secondary btn-sm button-margin" data-button="duplicate-chain"' +
              (permChainDublication ? '' : 'disabled')+'>Duplicate Chain</button>\t' +
          '<button id="edit-chain-' + id + '" class="btn btn-secondary btn-sm button-margin" data-button="edit-chain"' +
          '>Edit</button>\n' +
          '<button id="delete-chain-' + id + '" class="btn btn-secondary btn-sm button-margin" data-button="delete-chain" ' +
          'hidden>Delete</button>'
        }
      }
    ],
    'autoWidth': false,
    'createdRow': (row, data, index) => {
      $(row).attr('id', 'chain-row-' + data['ChainProfileID'])
    },
    'columnDefs': [{
      'targets': '_all',
      'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
        $(td).attr('id', 'chain-cell-' + colIndex + "-" + rowData['ChainProfileID'])
      }
    }],
    "drawCallback": ( settings ) => {
      setPageIds()
    },
  initComplete: function() {
      let input = $('#chainsTable_filter input').attr('id', 'chains-searchTerm').unbind();
      let self = this.api();

      // User must click this button to search now
      let $searchButton = $('<button>')
          .text('Search')
          .attr('id', 'chainSearchButton')
          .attr('class', 'btn btn-primary button-margin')
          .click(function() {
            self.search(input.val()).draw();
          });

      $('#chainsTable_filter').append($searchButton);

      // Add Event trigger to execute search when enter key is pressed
      addSearchKeyBind('chains-searchTerm', 'chainSearchButton');
    }
  })

  $('#chainsTable tbody').on('click', 'button', function () {

    if ($(this).attr('data-button') !== undefined ){
      let row = table.row($(this).parents('tr'))
      let data = row.data()
      const profID = data['ChainProfileID']
      const acquirerName = data['AcquirerName']

      switch ($(this).attr('data-button')) {
        case 'edit-chain':
          location.href = '/profileMaintenance?profileId=' + profID + '&type=chain'
          break
        case 'duplicate-chain':
          var tr = $(this).closest('tr');
          if (row.child.isShown()) {
            row.child.hide();
            tr.removeClass('shown');
          } else {
            row.child(format(profID,acquirerName)).show();
            tr.addClass('shown');
          }
          break
      }
    }
  })

}

function buildAcquirersTable (token) {
  let table = $('#acquirersTable').DataTable({
    ajax: {
      'url': location.href,
      'type': 'POST',
      'data': (d) => {
        return $.extend({}, d, {
          'requestType': 'acquirer',
          csrfmiddlewaretoken: token
        })
      }
    },
    deferRender: true,
    'processing': true,
    "deferLoading": 0,
    'serverSide': true,
    'pageLength': 15,
    'lengthMenu': [
      [10, 15, 25, 50, 100, 500, -1],
      [10, 15, 25, 50, 100, 500, 'All']
    ],
    'columns': [
      { data: 'AcquirerProfileID' },
      { data: 'AcquirerName' },
      {
        data: null,
        'orderable': false,
        'render': (cellData, cellType, rowData, meta) => {
          const id = rowData['AcquirerProfileID']
          return '<button id="edit-acquirer-' + id + '" class="btn btn-secondary btn-sm button-margin" data-button="edit-acquirer"' +
            '>Edit</button>\n' +
            '<button id="delete-acquirer-' + id + '" class="btn btn-secondary btn-sm button-margin" data-button="delete-acquirer" ' +
            'hidden>Delete</button>'
        }
      }
    ],
    'autoWidth': false,
    'createdRow': (row, data, index) => {
      $(row).attr('id', 'acquirer-row-' + data['AcquirerProfileID'])
    },
    'columnDefs': [{
      'targets': '_all',
      'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
        $(td).attr('id', 'acquirer-cell-' + colIndex + "-" + rowData['AcquirerProfileID'])
      }
    }],
    "drawCallback": ( settings ) => {
      setPageIds()
    },
    initComplete: function() {
      let input = $('#acquirersTable_filter input').attr('id', 'acquirers-searchTerm').unbind();
      let self = this.api();

      // User must click this button to search now
      let $searchButton = $('<button>')
          .text('Search')
          .attr('id', 'acquirerSearchButton')
          .attr('class', 'btn btn-primary button-margin')
          .click(function() {
            self.search(input.val()).draw();
          });

      $('#acquirersTable_filter').append($searchButton);

      // Add Event trigger to execute search when enter key is pressed
      addSearchKeyBind('acquirers-searchTerm', 'acquirerSearchButton');
    }
  });

  $('#acquirersTable tbody').on('click', 'button', function () {
    let row = table.row($(this).parents('tr'))
    let data = row.data()
    const profID = data['AcquirerProfileID']

    switch ($(this).attr('data-button')) {
      case 'edit-acquirer':
        location.href = '/profileMaintenance?profileId=' + profID + '&type=acquirer'
        break
    }
  })
}

function clearFilter () {
  $('#searchTerm')[0].value = ''
  search()
}

function setActiveTab (tab) {
  ActiveTab = tab
}

function exportSearch (filtered) {
  let exportModal = $('#exportModal')
  exportModal.modal({ backdrop: 'static', keyboard: false })
  exportModal.modal('show')
  const searchTerm = $('.tab-content div.active div[id$=_filter] input').val()

  $.ajax({
    url: 'exportSearch',
    contentType: 'application/x-www-form-urlencoded',
    data: {
      'SearchTerm': searchTerm,
      'ActiveTab': ActiveTab,
      'Filtered':filtered
    },
    success: (d, status, req) => {
      let clean = DOMPurify.sanitize(req.getResponseHeader('fileName'));
      window.location.href = "/downloadExportedReport/" + clean;
      $('#exportModal').modal('hide');
    },
    error: (d) => {
      $('#exportModal').modal('hide')
      setTimeout(function () {
        alert(d.responseText)
      }, 100);
    }
  })
}

function cancelExport () {
  const data = {
    csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
  }
  $.ajax({
    data: data,
    url: 'cancelExport',
    method: 'POST',
    success: function (d) {
      req.abort()
      $('#exportModal').modal('hide')
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function setPageIds () {
  let pageLinks = $('a[class="page-link"]');
  $(pageLinks).each(function(index,element){
    let content = $(this).html();
    // Only add this to page number links, not to actions like next or previous
    if($.isNumeric( content )) {
      $(this).attr('id', $(this).attr('aria-controls') + '-page-' + content);
    }
  });
}