'use strict'

let After = ''
let Before = ''
let Acquirer = ''
let Name = ''
let User = ''
let Module = ''

let userAuditOffset = 0

let auditFilterData = {
  Filters: []
}

$(document).ready(function () {
  Select()
  applyMultiSelects()
})

function applyMultiSelects () {
  const maxHeight = 200
  const numDisplayed = 1
  const buttonTemplate = '<button type="button" class="multiselect dropdown-toggle btn btn-outline-success" ' +
      'data-bs-toggle="dropdown">' +
      '<span class="multiselect-selected-text "></span>' +
      '</button>'

  $('#user').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'user-dropdown')
    }
  })

  $('#groups').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'groups-dropdown')
    }
  })

  $('#group').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'group-dropdown')
    }
  })

  $('#permissions').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'permissions-dropdown')
    }
  })

  $('#acquirers').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'acquirers-dropdown')
    }
  })

  $('#um-user-groups').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'um-user-groups-dropdown')
    }
  })

  $('#um-users').multiselect({
    maxHeight: maxHeight,
    numberDisplayed: numDisplayed,
    templates: {
      button: buttonTemplate,
    },
    onInitialized: (select) => {
      select.next().children('button').attr('id', 'um-users-dropdown')
    }
  })
}

function SaveUserGroup () {
  $.ajax({
    data: {
      User: $('#user')[0].value,
      Group: $('#group')[0].value,
      Groups: $('#groups').val(),
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    },
    type: 'POST',
    url: '/userManagement/SaveUserGroup',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementPartial').html(data)
      applyMultiSelects()
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function DeleteGroup () {
  confirmDialog('Delete Group?', 'Are you sure you want to delete Group? <br /><br /> Warning: <br /> If this group is assigned to existing users the permissions will be removed',
    function () {
      $.ajax({
        data: {
          User: $('#user')[0].value,
          Group: $('#group')[0].value,
          csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value,
        },
        type: 'POST',
        url: '/userManagement/DeleteGroup',
        success: function (data) {
          data = sanitizeHTML(data)
          $('#userManagementPartial').html(data)
          applyMultiSelects()
          Select()
        },
        error: function (data) {
          displayWarningMessage(data.responseText)
        }
      })
    })
}

function AddGroup () {
  hideWarning()
  $.ajax({
    data: {
      Name: $('#alterGroup')[0].value,
      User: $('#user')[0].value,
      Group: $('#group')[0].value,
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value,
    },
    type: 'POST',
    url: '/userManagement/AddGroup',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementPartial').html(data)
      applyMultiSelects()
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function RenameGroup () {
  confirmDialog('Rename Group?', 'Are you sure you want to rename Group?',
    function () {
      hideWarning()
      $.ajax({
        data: {
          Name: $('#alterGroup')[0].value,
          User: $('#user')[0].value,
          Group: $('#group')[0].value,
          csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value,
        },
        type: 'POST',
        url: '/userManagement/RenameGroup',
        success: function (data) {
          data = sanitizeHTML(data)
          $('#userManagementPartial').html(data)
          applyMultiSelects()
          Select()
        },
        error: function (data) {
          displayWarningMessage(data.responseText)
        }
      })
    })
}

function SaveGroupPermissions () {
  $.ajax({
    data: {
      User: $('#user')[0].value,
      Group: $('#group')[0].value,
      Permissions: $('#permissions').val(),
      AcquirerIDs: $('#acquirers').val(),
      AcquirerNames: $('#acquirers option:selected').toArray().map((a) => {return a.text}),
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value,
    },
    type: 'POST',
    url: '/userManagement/SaveGroupPermissions',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementPartial').html(data)
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function Select () {
  $.ajax({
    data: {
      User: $('#user')[0].value,
      Group: $('#group')[0].value,
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value,
    },
    type: 'POST',
    url: '/userManagement/Select',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementPartial').html(data)
      applyMultiSelects()
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function addUser () {
  $.ajax({
    data: {
      'username': $('#um-username')[0].value,
      'user-groups': JSON.stringify($('#um-user-groups').val()),
      'csrfmiddlewaretoken': $('input[name=csrfmiddlewaretoken]')[0].value,
    },
    type: 'POST',
    url: '/userManagement/addUser',
    success: (data) => {
      displayWarningMessage(data)
    },
    error: (data) => {
      displayWarningMessage(data.responseText)
    }
  })
}

function deleteUser () {
  $.ajax({
    data: {
      'username': $('#um-users').val(),
      'csrfmiddlewaretoken': $('input[name=csrfmiddlewaretoken]')[0].value,
    },
    type: 'POST',
    url: '/userManagement/deleteUser',
    success: (data) => {
      displayWarningMessage(data)
    },
    error: (data) => {
      displayWarningMessage(data.responseText)
    }
  })
}

//Retrieves the user change audit trail
function userHistoryTabSelect (tabIdent, offset) {
  $.ajax({
    data: {
      Filters: JSON.stringify(auditFilterData),
      Offset: offset,
      Type: tabIdent,
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    },
    type: 'POST',
    url: '/userChangeAuditHistory',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementAudit').html(data)
      $('#shown-results-' + tabIdent).html('Showing Results ' + offset + ' to ' + (offset + 50))
      $('#shown-results-footer-' + tabIdent).html('Showing Results ' + offset + ' to ' + (offset + 50))
      let prevResults = $('#prev-results-' + tabIdent)
      let prevResultsFooter = $('#prev-results-footer-' + tabIdent)
      if (offset > 0) {
        prevResults.removeClass('opacity')
        prevResults.addClass('clickable')
        prevResultsFooter.removeClass('opacity')
        prevResultsFooter.addClass('clickable')
      } else {
        prevResultsFooter.removeClass('clickable')
        prevResultsFooter.addClass('opacity')
      }
      initFilter()
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function initFilter () {
  $('#after-input').datetimepicker({ step: 5 })
  $('#before-input').datetimepicker({ step: 5 })
  $('#clear-filters').click(function () {
    clearFilters()
    console.log('Clearing')
  })
}

function clearFilters () {
  $('#after-input')[0].value = ''
  $('#before-input')[0].value = ''
  $('#acquirer-input')[0].value = ''
  $('#user-input')[0].value = ''
  $('#name-input')[0].value = ''
  $('#module-input')[0].value = ''
  auditFilterData.Filters = []
}

function setFilters () {
  After = $('#after-input')[0].value
  Before = $('#before-input')[0].value
  Acquirer = $('#acquirer-input')[0].value
  Name = $('#name-input')[0].value
  User = $('#user-input')[0].value
  Module = $('#module-input')[0].value

  auditFilterData.Filters.push({
    Name: 'after',
    Value: After
  })
  auditFilterData.Filters.push({
    Name: 'before',
    Value: Before
  })
  auditFilterData.Filters.push({
    Name: 'acquirer',
    Value: Acquirer
  })
  auditFilterData.Filters.push({
    Name: 'name',
    Value: Name
  })
  auditFilterData.Filters.push({
    Name: 'updatedby',
    Value: User
  })
  auditFilterData.Filters.push({
    Name: 'module',
    Value: Module
  })
}

function filterAuditHistory (tabIdent, offset) {
  auditFilterData.Filters = []
  setFilters()

  $.ajax({
    data: {
      Filters: JSON.stringify(auditFilterData),
      Offset: offset,
      Type: tabIdent,
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    },
    type: 'POST',
    url: '/userChangeAuditHistory',
    success: function (data) {
      data = sanitizeHTML(data)
      $('#userManagementAudit').html(data)
      $('#shown-results-' + tabIdent).html('Showing Results ' + offset + ' to ' + (offset + 50))
      $('#shown-results-footer' + tabIdent).html('Showing Results ' + offset + ' to ' + (offset + 50))
      let prevResults = $('#prev-results-' + tabIdent)
      let prevResultsFooter = $('#prev-results-footer-' + tabIdent)
      if (offset > 0) {
        prevResults.removeClass('opacity')
        prevResults.addClass('clickable')
        prevResultsFooter.removeClass('opacity')
        prevResultsFooter.addClass('clickable')
      } else {
        prevResults.removeClass('clickable')
        prevResults.addClass('opacity')
        prevResultsFooter.removeClass('clickable')
        prevResultsFooter.addClass('opacity')
      }
      initFilter()

      $('#after-input')[0].value = After
      $('#before-input')[0].value = Before
      $('#acquirer-input')[0].value = Acquirer
      $('#name-input')[0].value = Name
      $('#user-input')[0].value = User
      $('#module-input')[0].value = Module
    },
    error: function (data) {
      displayWarningMessage(data.responseText)
    }
  })
}

function checkAuditKeyPress () {
  let key = event.which
  if (key === 13) {
    filterAuditHistory('Item', 0)
  }
}

function userAuditPageResults (pageAmount, tabIdent) {
  if (tabIdent === 'userAudit' && (pageAmount >= 0 || userAuditOffset > -1)) {
    userAuditOffset = userAuditOffset + pageAmount
    filterAuditHistory(tabIdent, userAuditOffset)
  }
}

function exportAuditHistory () {
  auditFilterData.Filters = []
  setFilters()
  let filters = 'Filters=' + JSON.stringify(auditFilterData) + '&' + 'Offset=' + userAuditOffset

  req.open('Get', '/exportUserChangeAuditHistory?' + filters, true)
  req.setRequestHeader('Content-type', 'application/x-www-form-urlencoded')

  req.responseType = 'blob'
  req.onload = function (event) {
    var blob = req.response
    var fileName = req.getResponseHeader('fileName') //if you have the fileName header available
    var link = document.createElement('a')
    link.href = encodeURI(window.URL.createObjectURL(blob))
    link.download = fileName
    link.click()
  }
  req.send(filters)
}

