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
    case "sitemap":
      url = "/sitemap";
      formData.append("url", inputFieldValue);
      break;
    case "crawlcheck":
      url = "/crawlcheck";
      formData.append("siteURL", inputFieldValue);
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
            case "SSLInfo":
            case "cookie":
            case "whois":
            case "serverInfo":
            case "portScanner":
            case "sitemap":
            case "crawlcheck":
              responseContainer.innerHTML = formatData(jsonData);
              break;
            case "screenshot":
              responseContainer.innerHTML = `<img src="data:image/png;base64,${jsonData.ScreenshotBase64}" alt="Screenshot">`;
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
        const value = data[key];
        formatted += `<span style="color: aqua; font-weight: bold">${key}:</span> <span style="color: #00FF00">${formatData(value)}</span><br>`;
      }
    }
  } else {
    formatted += data;
  }

  formatted += '</div>';
  return formatted;
}


