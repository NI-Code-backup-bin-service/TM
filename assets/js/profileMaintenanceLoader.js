 $(document).ready(() => {
  $.ajax({
    url: '/profileMaintenance?type='+Type,
    data: {
      profileID: ID,
      DGUpdated: DGUpdated,
      csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    },
    method: 'POST',
    success: function (data) {
      $("body").append(data)
      $("#profLoader").hide()
    },
    error: function (data) {
      if (data.status === 401 || data.status === 403) {
        displayWarningMessage('User not authorised')
      } else {
        displayWarningMessage(data.responseText)
      }
    }
  })
})