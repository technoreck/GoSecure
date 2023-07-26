document.getElementById("functionalityForm").addEventListener("submit", function (event) {
  event.preventDefault(); 

  const selectedFunctionality = document.querySelector('input[name="functionality"]:checked').value;
  const inputFieldValue = document.getElementById("inputField").value;

  let url;
  let formData = new FormData();

  switch (selectedFunctionality) {
    case "dnsLookup":
      url = "/dnsinfo";
      formData.append("hostname", inputFieldValue);
      break;
    case "getData":
      url = `/getData?url=${encodeURIComponent(inputFieldValue)}`;
      break;
    case "hstsChecker":
      url = "/hsts";
      formData.append("url", inputFieldValue);
      break;
    case "portScanner":
      url = "/scan";
      formData.append("hostname", inputFieldValue);
      break;
    case "servstat":
      url = "/servs";
      formData.append("url", inputFieldValue);
      break;
    case "dns":
      url = "/resolve";
      formData.append("url", inputFieldValue);
      break;
    case "dnssec":
      url = "/dnssec";
      formData.append("url", inputFieldValue);
      break;
    case "screenshot":
      url = "/screenshot";
      formData.append("url", inputFieldValue);
      break;
    default:
      alert("Please select a functionality.");
      return;
  }

  if (selectedFunctionality === "getData") {
    fetch(url)
      .then(response => {
        if (!response.ok) {
          throw new Error("Server responded with an error status.");
        }
        return response.json();
      })
      .then(data => {
        const responseContainer = document.getElementById("responseContainer");
        responseContainer.innerHTML = formatData(data);
      })
      .catch(error => {
        const responseContainer = document.getElementById("responseContainer");
        responseContainer.textContent = "An error occurred: " + error.message;
      });
  } else {
    fetch(url, {
      method: "POST",
      body: formData,
    })
      .then(response => {
        if (!response.ok) {
          throw new Error("Server responded with an error status.");
        }
        return response.text();
      })
      .then(data => {
        const responseContainer = document.getElementById("responseContainer");

        try {
          const jsonData = JSON.parse(data);

          switch (selectedFunctionality) {
            case "dnsLookup":
            case "hstsChecker":
            case "servstat":
            case "dns":
            case "dnssec":
              responseContainer.innerHTML = formatData(jsonData);
              break;
            case "portScanner":
              responseContainer.innerHTML = formatPortScannerData(jsonData);
              break;
            case "screenshot":
              responseContainer.innerHTML = `<img src="data:image/png;base64,${jsonData.ScreenshotBase64}" alt="Screenshot">`;
              break;
            case "serverInfo":
              responseContainer.innerHTML = formatServerInfoData(jsonData);
              break;
            default:
              responseContainer.innerHTML = "Invalid functionality.";
          }
        } catch (error) {
          responseContainer.innerHTML = `<div class="data-item">${data}</div>`;
        }
      })
      .catch(error => {
        const responseContainer = document.getElementById("responseContainer");
        responseContainer.textContent = "An error occurred: " + error.message;
      });
  }
});

function formatData(data) {
  let formatted = '<div class="data-item">';
  
  if (typeof data === 'object') {
    if (Array.isArray(data)) {
      formatted += '<strong>Data Items:</strong><br>';
      data.forEach(item => {
        formatted += formatData(item) + '<br>';
      });
    } else {
      for (let key in data) {
        formatted += `<strong>${key}:</strong> ${formatData(data[key])}<br>`;
      }
    }
  } else {
    formatted += data;
  }

  formatted += '</div>';
  return formatted;
}

function formatPortScannerData(data) {
  let formatted = '<div class="data-item"><strong>Port Scanner Results:</strong><br>';
  formatted += JSON.stringify(data);
  formatted += '</div>';
  return formatted;
}

function formatServerInfoData(data) {
  let formatted = '<div class="data-item"><strong>Server Info:</strong><br>';

  for (let key in data) {
    let value = data[key];

    if (Array.isArray(value)) {
      formatted += `<strong>${key}:</strong><br>`;
      value.forEach(item => {
        if (typeof item === 'object') {
          formatted += '<ul>';
          for (let subKey in item) {
            formatted += `<li><strong>${subKey}:</strong> ${item[subKey]}</li>`;
          }
          formatted += '</ul>';
        } else {
          formatted += `${JSON.stringify(item)}<br>`;
        }
      });
    } else if (typeof value === 'object') {
      formatted += `<strong>${key}:</strong><br>`;
      formatted += `${JSON.stringify(value)}<br>`;
    } else {
      formatted += `<strong>${key}:</strong> ${value}<br>`;
    }
  }

  formatted += '</div>';
  return formatted;
}
