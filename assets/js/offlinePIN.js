'use strict'

const op = 'op-'

const serial = 'serial'
const imei = 'imei'
const offlineExpiry = 'expiry1'
const resetExpiry = 'expiry2'

const snInputPrefix = op + serial + '-'
const imeiInputPrefix = op + imei + '-'
const offlineExpiryInputPrefix = op + offlineExpiry + '-'
const resetExpiryInputPrefix = op + resetExpiry + '-'

const offlinePINSN = 'offline-sn'
const offlinePINIMEI = 'offline-imei'
const resetPINSN = 'reset-sn'
const resetPINIMEI = 'reset-imei'
const offlineExpiresAfter = 'expiry-date-engineering'
const resetExpiresAfter = 'expiry-date-reset'
const deleteRowButton = 'delete-row-button'

const offlinePINSNPrefix = op + offlinePINSN + '-'
const offlinePINIMEIPrefix = op + offlinePINIMEI + '-'
const resetPINSNPrefix = op + resetPINSN + '-'
const resetPINIMEIPrefix = op + resetPINIMEI + '-'
const offlineExpiresAfterPrefix = op + offlineExpiresAfter + '-'
const resetExpiresAfterPrefix = op + resetExpiresAfter + '-'
const deleteRowButtonPrefix = op + deleteRowButton + '-'

$(document).ready(() => {
  const table = setupTable()

  table.on('click', '.op-delete-row', function () {
    if (table.rows().count() === 1) {
      // Don't allow removal of single row
      return
    }

    table
      .row($(this).parents('tr'))
      .remove()
      .draw()
  });

  $('#op-add-row').click(() => {
    addRow(table)
  })
  // Add single empty row
  addRow(table)

  setupRemoveAllRows(table)

  setupGeneratePIN(table)

  setupCSVImport()
})

function setupTable () {
  const table = $('#op-table').DataTable({
    'dom': 'Bfrtip',
    'deferRender': true,
    'searching': false,
    'ordering': false,
    'info': false,
    'autoWidth': false,
    'columnDefs': [
      { 'width': '13%', 'targets': [0, 1, 2, 3] },
      { 'width': '8%', 'targets': [4, 5, 6, 7, 8, 9] }
    ],
    'buttons': [
      {
        'extend': 'csv',
        'filename': generateFilename,
        'text': 'Export CSV',
        'className': 'btn-success',
        'action': function (e, dt, node, config) {
          // overriding the action is necessary to display the spinner
          const spinner = $('#op-loader')
          spinner.toggle(true)

          const button = $(node)
          button.attr('disabled', 'disabled')

          setTimeout(() => {
            // this is the default, built-in, csv export for datatables buttons
            $.fn.dataTable.ext.buttons.csvHtml5.action.call(this, e, dt, node, config)

            spinner.toggle(false)
            button.removeAttr('disabled')
          }, 50)
        },
        'exportOptions': {
          'columns': [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
          'format': {
            'body': function (html, row, col, node) {
              if (col > 3) { // if not first 4 rows (input rows), do nothing
                // return default formatting
                return $.fn.DataTable.Buttons.stripData(html, null)
              }

              const inputCells = $(node).find('input')
              // there should only be 1, so just grab the first
              return inputCells[0].value
            }
          }
        },
        'fieldBoundary': ''
      }
    ]
  })

  // place the export csv button in the correct place
  table.buttons().container().appendTo($('#op-csv-export-container'))

  return table
}

// rowIndex is used so that each new row is created with a new unique index
let rowIndex = 0

// serial, imei and expiry use default values for add row button
function addRow (table, serial = '', imei = '', offlineExpiry = 1, resetExpiry = 1) {
  const index = rowIndex++

  const snInputID = snInputPrefix + index
  const snInput = $('<input class="input-full-width">')
  snInput.attr({
    type: 'text',
    id: snInputID,
    name: snInputID,
    value: serial
  })
  const snColumn = $('<td>').append(snInput)

  const imeiInputID = imeiInputPrefix + index
  const imeiInput = $('<input class="input-full-width">')
  imeiInput.attr({
    type: 'text',
    id: imeiInputID,
    name: imeiInputID,
    value: imei
  })
  const imeiColumn = $('<td>').append(imeiInput)

  const offlineExpiryInputID = offlineExpiryInputPrefix + index
  const offlineExpiryInput = $('<input class="input-full-width" min="1" max="5">')
  offlineExpiryInput.change(limitInputValue)
  offlineExpiryInput.attr({
    type: 'number',
    id: offlineExpiryInputID,
    name: offlineExpiryInputID,
    value: offlineExpiry
  })
  const offlineExpiryColumn = $('<td>').append(offlineExpiryInput)

  const resetExpiryInputID = resetExpiryInputPrefix + index
  const resetExpiryInput = $('<input class="input-full-width" min="1" max="5">')
  resetExpiryInput.change(limitInputValue)
  resetExpiryInput.attr({
    type: 'number',
    id: resetExpiryInputID,
    name: resetExpiryInputID,
    value: resetExpiry
  })
  const resetExpiryColumn = $('<td>').append(resetExpiryInput)

  const offlinePINSNCell = $('<td>')
  offlinePINSNCell.attr('id', offlinePINSNPrefix + index)

  const offlinePINIMEICell = $('<td>')
  offlinePINIMEICell.attr('id', offlinePINIMEIPrefix + index)

  const offlineExpiresAfterCell = $('<td>')
  offlineExpiresAfterCell.attr('id', offlineExpiresAfterPrefix + index)

  const resetPINSNCell = $('<td>')
  resetPINSNCell.attr('id', resetPINSNPrefix + index)

  const resetPINIMEICell = $('<td>')
  resetPINIMEICell.attr('id', resetPINIMEIPrefix + index)

  const resetExpiresAfterCell = $('<td>')
  resetExpiresAfterCell.attr('id', resetExpiresAfterPrefix + index)

  const deleteRowButton = $('<button class="btn btn-secondary btn-sm op-delete-row">Delete</button>')
  deleteRowButton.attr('id', deleteRowButtonPrefix + index)
  const deleteRowCell = $('<td>').append(deleteRowButton)

  const row = $('<tr>')
    .append(snColumn, imeiColumn, offlineExpiryColumn, resetExpiryColumn, offlinePINSNCell, offlinePINIMEICell,
      offlineExpiresAfterCell, resetPINSNCell, resetPINIMEICell, resetExpiresAfterCell, deleteRowCell)

  table.row.add(row).draw(false)
}

function limitInputValue (e) {
  e.target.value = limitDate(e.target.value)
}

function limitDate (days) {
  return Math.min(Math.max(days, 1), 5)
}

function sanitiseNumber (days) {
  const number = Number(days)
  if (Number.isInteger(number)) {
    return Math.floor(number)
  }

  return 1
}

function displayPINGenerationLoader () {
  $('#op-loader').toggle(true)

  $('#op-gen-pin').attr('disabled', 'disabled')
  $('#op-gen-pin-engineering').attr('disabled', 'disabled')
  $('#op-gen-pin-reset').attr('disabled', 'disabled')
}

function hidePINGenerationLoader () {
  $('#op-gen-pin').removeAttr('disabled')
  $('#op-gen-pin-engineering').removeAttr('disabled')
  $('#op-gen-pin-reset').removeAttr('disabled')

  $('#op-loader').toggle(false)
}

function getTableData (table) {
  let data = []
  table.rows().every(function (index, tableLoop, rowLoop) {
    const cells = $(this.node()).find('input')
    let rowData = {}
    for (const cell of cells) {
      rowData[cell.name.split('-')[1]] = cell.value
    }
    if (rowData[serial] === '' && rowData[imei] === '') {
      // invalid row
      return
    }
    data.push(rowData)
  })

  if (table.rows().count() !== data.length) {
    // if there are less data entries than number of rows, not all rows are valid
    displayWarningMessage('Please include either a Serial Number or IMEI for every entry')
    return null
  }
  return data
}

function performPINGeneration (table, mode = 'both') {
  displayPINGenerationLoader()

  const data = getTableData(table)
  if (data === null) {
    hidePINGenerationLoader()
    return
  }

  $.ajax({
    url: '/generateOfflinePIN',
    data: {
      pedData: JSON.stringify(data),
      mode: mode,
      csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
    },
    method: 'POST',
    success: (data) => {
      const results = JSON.parse(data)
      writeResultsToTable(results)
    },
    error: (data) => {
      hidePINGenerationLoader()
      displayError(data)
    }
  })
}

function setupRemoveAllRows (table) {
  $('#op-clear-rows').click(() => {
    table.clear().draw()
    addRow(table)
  })
}

function setupGeneratePIN (table) {
  $('#op-gen-pin').click(() => {
    performPINGeneration(table)
    return false
  })

  $('#op-gen-pin-engineering').click(() => {
    performPINGeneration(table, 'offline')
    return false
  })

  $('#op-gen-pin-reset').click(() => {
    performPINGeneration(table, 'reset')
    return false
  })
}

function setupCSVImport () {
  $('#op-csv-upload').submit((e) => {
    e.preventDefault()

    const spinner = $('#op-loader')
    spinner.toggle(true)

    const button = $('#op-csv-upload-button')
    button.attr('disabled', 'disabled')

    // setTimeout is used here so previous spinner toggle and button attributes are updated before page processing happens
    setTimeout(() => {
      const formData = new FormData(document.querySelector('#op-csv-upload'))

      $.ajax({
        url: '/offlinePINImportCSV',
        data: formData,
        method: 'POST',
        processData: false,
        contentType: false,
        success: (data) => {
          button.removeAttr('disabled')
          spinner.toggle(false)
          displayWarningMessage('CSV data imported')
          const results = JSON.parse(data)
          appendDataToTable(results)
        },
        error: (data) => {
          button.removeAttr('disabled')
          displayError(data)
        }
      })
    }, 50)
  })
}

function displayError (data) {
  const spinner = $('#op-loader')
  spinner.toggle(false)
  if (data.status === 401 || data.status === 403 ){
    displayWarningMessage('User not authorised')
  } else {
    displayWarningMessage(data.responseText)
  }
}

function writeResultsToTable (data) {
  const table = $('#op-table').DataTable()

  table.rows().every(function (index, tableLoop, rowLoop) {
    setTimeout(() => {
      displayWarningMessage('Processing entry ' + (index + 1) + ' out of ' + data.length)

      const inputCells = $(this.node()).find('input')

      const rowData = data[index]

      let rowIndex = -1
      // Below loop is for data validation
      for (const cell of inputCells) {
        const split = cell.name.split('-')
        const name = split[1]
        rowIndex = split[2]

        const value = cell.value
        const returnedValue = rowData[name]

        // This check ensures that we never display the PIN for another PED if the table rows get mismatched
        if (returnedValue != null && returnedValue !== value) {
          hidePINGenerationLoader()
          displayWarningMessage('Returned data does not match sent data')
          return
        }
      }

      if (rowIndex === -1) {
        hidePINGenerationLoader()
        displayWarningMessage('Invalid row index found in returned data')
        return
      }

      this.cell(index, 4).data(rowData[offlinePINSN])
      this.cell(index, 5).data(rowData[offlinePINIMEI])
      this.cell(index, 6).data(rowData[offlineExpiresAfter])
      this.cell(index, 7).data(rowData[resetPINSN])
      this.cell(index, 8).data(rowData[resetPINIMEI])
      this.cell(index, 9).data(rowData[resetExpiresAfter])

      if (index+1 === data.length) {
        setTimeout(() => {
          hidePINGenerationLoader()
          displayWarningMessage('PIN(s) successfully generated')
        }, 50)
      }
    }, 0)
  })

  table.draw()
}

function appendDataToTable (data) {
  const table = $('#op-table').DataTable()

  // Remove last row if empty, this also covers the default row
  const lastRow = table.row( ':last')
  const cells = $(lastRow.node()).find('input')
  let lastEmpty = true
  for (const cell of cells) {
    const cellName = cell.name.split('-')[1]
    if ((cellName === serial || cellName === imei) && cell.value !== '') {
      // if serial number or IMEI cell is not empty, don't delete last row
      lastEmpty = false
      break
    }
  }
  if (lastEmpty) {
    lastRow.remove()
  }

  for (const row of data) {
    const offlineExpiryDays = limitDate(sanitiseNumber(row[offlineExpiry]))
    const resetExpiryDays = limitDate(sanitiseNumber(row[resetExpiry]))

    addRow(table, row[serial], row[imei], offlineExpiryDays, resetExpiryDays)
  }
}

function generateFilename () {
  const time = new Date().toISOString()
  const dtsplit = time.split('T')
  const timesplit = dtsplit[1].split('.')
  const filename = dtsplit[0].replace(/-/g,'') + '-' + timesplit[0].replace(/:/g,'') + '.csv'
  return filename
}