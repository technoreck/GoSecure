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
    case "SSLInfo":
      url = "/sslinfo";
      formData.append("url", inputFieldValue);
      break;
    case "cookie":
      url = "/cookie";
      formData.append("url", inputFieldValue);
      break;
    case "whois":
      url = "/whois";
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
  }
  else if (selectedFunctionality === 'SSLInfo') {
    fetch(url, {
      method: "POST",
      body: formData,
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Server responded with an error status.');
        }
        return response.json();
      })
      .then(data => {
        const responseContainer = document.getElementById('responseContainer');

        if (data.error) {
          // Display error message if there's an error in the response
          responseContainer.innerHTML = `<p>Error: ${data.error}</p>`;
        } else if (data.certificates && data.certificates.length > 0) {
          // Format the SSL information and display it
          let formatted = '<div class="data-item">';
          data.certificates.forEach(cert => {
            formatted += `<strong>Subject:</strong> ${cert.subject}<br>`;
            formatted += `<strong>Issuer:</strong> ${cert.issuer}<br>`;
            formatted += `<strong>Valid From:</strong> ${cert.valid_from}<br>`;
            formatted += `<strong>Valid Until:</strong> ${cert.valid_until}<br>`;
            formatted += `<strong>Serial Number:</strong> ${cert.serial_number}<br>`;
            formatted += `<strong>Signature Algorithm:</strong> ${cert.signature_algorithm}<br>`;
            formatted += `<strong>Key Usage:</strong> ${cert.key_usage}<br>`;
            formatted += `<strong>Is CA Cert:</strong> ${cert.is_ca_cert}<br>`;
            if (cert.dns_names && cert.dns_names.length > 0) {
              formatted += `<strong>DNS Names:</strong> ${cert.dns_names.join(', ')}<br>`;
            }
            formatted += '<br>';
          });
          formatted += '</div>';
          responseContainer.innerHTML = formatted;
        } else {
          responseContainer.innerHTML = '<p>No SSL information found.</p>';
        }
      })
      .catch(error => {
        const responseContainer = document.getElementById('responseContainer');
        responseContainer.innerHTML = 'An error occurred: ' + error.message;
      });
  }
  else if (selectedFunctionality === 'cookie') {
    fetch(url, {
      method: "POST",
      body: formData,
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Server responded with an error status.');
        }
        return response.json();
      })
      .then(data => {
        const responseContainer = document.getElementById('responseContainer');
  
        if (data.error) {
          // Display error message if there's an error in the response
          responseContainer.innerHTML = `<p>Error: ${data.error}</p>`;
        } else if (data.cookies && data.cookies.length > 0) {
          // Format the cookie information and display it
          let formatted = '<div class="data-item">';
          data.cookies.forEach(cookie => {
            formatted += `<strong>Name:</strong> ${cookie.name}<br>`;
            formatted += `<strong>Value:</strong> ${cookie.value}<br>`;
            formatted += '<br>';
          });
          formatted += '</div>';
          responseContainer.innerHTML = formatted;
        } else {
          responseContainer.innerHTML = '<p>No cookies found.</p>';
        }
      })
      .catch(error => {
        const responseContainer = document.getElementById('responseContainer');
        responseContainer.innerHTML = 'An error occurred: ' + error.message;
      });
  }
  else if (selectedFunctionality === 'whois') {
    fetch(url, {
      method: "POST",
      body: formData,
    })
      .then(response => {
        if (!response.ok) {
          throw new Error("Server responded with an error status.");
        }
        return response.json(); // Parse the response as JSON
      })
      .then(data => {
        const responseContainer = document.getElementById("responseContainer");
        responseContainer.innerHTML = formatWHOISData(data);
      })
      .catch(error => {
        const responseContainer = document.getElementById("responseContainer");
        responseContainer.textContent = "An error occurred: " + error.message;
      });
  }
  else {
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

const form = document.getElementById('functionalityForm');
form.addEventListener('submit', handleFormSubmit);

function formatWHOISData(data) {
  // Create a formatted HTML representation of the WHOIS information
  let formatted = '<div class="data-item">';

  // Check if the response contains an "error" field
  if (data.error) {
    formatted += `<p>Error: ${data.error}</p>`;
  } else {
    // Iterate through the keys and values of the WHOIS data
    for (const key in data) {
      if (key !== 'termsOfUse' && key !== 'rawData') { // Exclude the "termsOfUse" and "rawData" fields
        const value = data[key];
        formatted += `<strong>${key}:</strong> ${value}<br>`;
      }
    }
  }

  formatted += '</div>';
  return formatted;
}
