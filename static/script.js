document.getElementById("dnsLookupForm").addEventListener("submit", function (event) {
    event.preventDefault();
    const hostname = document.getElementById("hostname").value;

    fetch("/dnsinfo", {
        method: "POST",
        body: new URLSearchParams({ hostname }),
    })
        .then(response => response.json())
        .then(data => {
            let result = "";

            for (const key in data) {
                if (Array.isArray(data[key])) {
                    const values = data[key].join("<br>");
                    result += `<div class="data-item"><b>${key}:</b><br>${values}</div>`;
                } else {
                    result += `<div class="data-item"><b>${key}:</b> ${data[key]}</div>`;
                }
            }

            document.getElementById("dnsResponse").innerHTML = result;
        })
        .catch(error => {
            document.getElementById("dnsResponse").textContent = "An error occurred: " + error.message;
        });
});

document.getElementById('getDataForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const url = document.getElementById('url').value;

    fetch(`/getData?url=${encodeURIComponent(url)}`)
        .then(response => response.json())
        .then(data => {
            let formattedData = formatData(data);
            document.getElementById('getDataResponse').innerHTML = formattedData;
        })
        .catch(error => {
            console.error(error);
            document.getElementById('getDataResponse').innerText = `Error: ${error.message}`;
        });
});

document.getElementById('portScannerForm').addEventListener('submit', async (event) => {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const response = await fetch('/scan', {
        method: 'POST',
        body: formData,
    });
    const result = await response.text();
    document.getElementById('result').innerText = result;
});

function formatData(data) {
    let formatted = '';
    for (let key in data) {
        let values = Array.isArray(data[key]) ? data[key] : [data[key]];
        formatted += `<div class="data-item"><strong>${key}:</strong><br>`;
        values.forEach(value => {
            formatted += `${value}<br>`;
        });
        formatted += '</div>';
    }
    return formatted;
}
document.getElementById('hstsForm').addEventListener('submit', async (event) => {
    event.preventDefault(); // Prevent the default form submission behavior.
    const form = event.target;
    const formData = new FormData(form);
    
    const response = await fetch('/hsts', {
        method: 'POST',
        body: formData,
    });
    
    if (response.ok) {
        try {
            // Try to parse the response as JSON.
            const responseData = await response.json();

            // Format the JSON data using the formatData function.
            const formattedOutput = formatData(responseData);

            // Display the formatted output in an element with the id 'result'.
            document.getElementById('hstsResult').innerHTML = formattedOutput;
        } catch (error) {
            // If parsing as JSON fails, display the response as plain text.
            const responseText = await response.text();
            document.getElementById('hstsResult').innerText = responseText;
        }
    } else {
        // Handle non-200 status code (e.g., error response from the server).
        document.getElementById('hstsResult').innerText = 'Error: Unable to process the request.';
    }
});

document.getElementById('servstat').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent the default form submission behavior.
    const url = document.getElementById('url').value;
    fetch('/servs', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'url=' + encodeURIComponent(url),
    })
    .then(response => response.text())
    .then(data => {
        document.getElementById('servstatResponse').innerHTML = data; // Display the result in the element with ID "servstatResponse"
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('servstatResponse').innerText = 'An error occurred during the scan.';
    });
});

document.getElementById('dnssecForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const urlInput = document.getElementById('dnssecUrl').value;
    
    fetch('/dnssec', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'url=' + encodeURIComponent(urlInput),
      })
        .then((response) => response.json())
        .then((data) => {
          const dnssecContainer = document.getElementById('dnssecResponse');
          dnssecContainer.innerHTML = '';
  
          if (data.error) {
            // If there is an error, display it in the response container
            dnssecContainer.innerHTML = 'Error fetching DNSSEC records: ' + data.error;
          } else {
            // If data is received, display RRSIG and DNSKEY records
            const rrsigRecords = data.RRIGRecords;
            const dnskeyRecords = data.DNSKEYRecords;
  
            // Display RRSIG records
            dnssecContainer.innerHTML += '<h3>RRSIG Records:</h3>';
            if (rrsigRecords.length > 0) {
              for (const rrsigRecord of rrsigRecords) {
                dnssecContainer.innerHTML += `<p>${JSON.stringify(rrsigRecord)}</p>`;
              }
            } else {
              dnssecContainer.innerHTML += '<p>No RRSIG records found.</p>';
            }
  
            // Display DNSKEY records
            dnssecContainer.innerHTML += '<h3>DNSKEY Records:</h3>';
            if (dnskeyRecords.length > 0) {
              for (const dnskeyRecord of dnskeyRecords) {
                dnssecContainer.innerHTML += `<p>${JSON.stringify(dnskeyRecord)}</p>`;
              }
            } else {
              dnssecContainer.innerHTML += '<p>No DNSKEY records found.</p>';
            }
          }
        })
        .catch((error) => {
          console.error('Error fetching DNSSEC records:', error);
          const dnssecContainer = document.getElementById('dnssecResponse');
          dnssecContainer.innerHTML = 'Error fetching DNSSEC records. Please try again later.';
      });
});

document.getElementById('screenshotForm').addEventListener('submit', function (event) {
    event.preventDefault();
    const urlInput = document.getElementById('screenshotUrl').value;
    
    fetch('/screenshot', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'url=' + encodeURIComponent(urlInput),
      })
        .then((response) => response.json())
        .then((data) => {
          const screenshotContainer = document.getElementById('SSResult');
          const img = new Image();
          img.src = 'data:image/png;base64,' + data.ScreenshotBase64;
          screenshotContainer.innerHTML = '';
          screenshotContainer.appendChild(img);
        })
        .catch((error) => {
          console.error('Error fetching screenshot:', error);
          const screenshotContainer = document.getElementById('SSResult');
          screenshotContainer.innerHTML = 'Error fetching screenshot. Please check the URL and try again.';
        });
});

function formatObject(obj, level = 0) {
    let formattedData = '';
    for (const key in obj) {
      const value = obj[key];
      if (typeof value === 'object' && value !== null) {
        formattedData += `${'\t'.repeat(level)}${key}\n${formatObject(value, level + 1)}`;
      } else {
        formattedData += `${'\t'.repeat(level)}${key}: ${value}\n`;
      }
    }
    return formattedData;
  }
// Example code for handling the "Server Info" response (dnsForm)
  const dnsForm = document.getElementById('dnsForm');
  const dnsResponseContainer = document.getElementById('dnsResponseContainer');
  
  dnsForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const formData = new FormData(dnsForm);
    const url = formData.get('url');
  
    try {
      const response = await fetch('/resolve', {
        method: 'POST',
        body: formData,
      });
      const data = await response.json();
  
      // Format the nested object data and display it in the response container
      const formattedData = formatObject(data);
      dnsResponseContainer.innerHTML = formattedData;
    } catch (error) {
      dnsResponseContainer.innerHTML = 'Error occurred while fetching data.';
    }
  });
