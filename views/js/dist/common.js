$('.hrefModal').click(function(){
    //alert("Modal Box hrefModal")
    var modalUrls=$(this).attr('data-href');
    var tid=$(this).attr('data-tid');

    //alert(modalUrls)

    $('#transModal').modal('show');
    $('#transModal').modal('show').find('.modal-body').load(modalUrls);
    $('#transModal .modal-dialog').css({"max-width":"80%", "margin-top": "20px"});
    $('#transModal .modal-title').html("Transaction Details - " + tid);
});

$(".loaderx").click(function(){
    //alert(1111)
$(".loaderx").html("<i class='fa-solid fa-spinner fa-spin-pulse'></i>");
});


// For Checkout Page JS

// Function to format time as MM:SS
function formatTime(seconds) {
    const minutes = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}

// The jQuery function from Step 1 goes here
function getUrlParameter(name) {
            name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
            var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
            var results = regex.exec(window.location.href);
            return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
        }

        $('#replybtn').click(function(){

          location.reload(true);
      });
    
      $("#crypto_network").change(function(){
        
        var myTid=$("#crypto_network option:selected").attr("data-id");
        $('#crypto_id').val(myTid);
    });


      function fetchStates() {
        //alert(111111111)
        const crypto_code = document.getElementById("crypto_code").value;
        const stateSelect = document.getElementById("crypto_network");

   // Get the dropdown element by its ID
    var cryptoDropdown = document.getElementById("crypto_code");
    // Get the selected option
    var selectedOption = cryptoDropdown.options[cryptoDropdown.selectedIndex];
    // Get the value of the custom attribute
    var dataCode = selectedOption.getAttribute("data-title");
    document.getElementById("crypto_title").value = dataCode;

        // Clear existing options
        stateSelect.innerHTML = '<option value="" >Select Network </option>';
    
        if (crypto_code) {
            fetch(`/get-network?crypto_code=${crypto_code}`)
            .then(response => response.json())
            .then(networks => {
                networks.forEach(network => {
                    const option = document.createElement("option");
                    option.value = network.Crypto_network;
                    option.setAttribute('data-id', network.Crypto_id)
                    option.textContent = network.Crypto_network;
                    stateSelect.appendChild(option);
                    
                });
            })
                .catch(error => console.error('Error fetching states:', error));
        }
    }


    function fetchnetwork() {
        //alert(111111111)
        const crypto_code = document.getElementById("crypto_code").value;
        const stateSelect = document.getElementById("crypto_network");

   // Get the dropdown element by its ID
    var cryptoDropdown = document.getElementById("crypto_code");
    // Get the selected option
    var selectedOption = cryptoDropdown.options[cryptoDropdown.selectedIndex];
    // Get the value of the custom attribute
    var dataCode = selectedOption.getAttribute("data-title");
    document.getElementById("crypto_title").value = dataCode;

        // Clear existing options
        stateSelect.innerHTML = '<option value="" >Select Network </option>';
    
        if (crypto_code) {
            fetch(`/get-network?crypto_code=${crypto_code}&tab=coin_list`)
            .then(response => response.json())
            .then(networks => {
                networks.forEach(network => {
                    const option = document.createElement("option");
                    option.value = network.Coin_network;
                    option.setAttribute('data-id', network.Coin_id)
                    option.textContent = network.Coin_network;
                    stateSelect.appendChild(option);
                    
                });
            })
                .catch(error => console.error('Error fetching states:', error));
        }
    }

//For change Date Format
    // Function to format the date
function formatDate(dateStr) {
    const date = new Date(dateStr);
    const options = {
        year: 'numeric',
        month: 'numeric', //long
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',  // Include seconds if needed
        hour12: true         // Set to false for 24-hour format
    };
    return date.toLocaleDateString(undefined, options);
}

// Get all date cells
const dateCells = document.querySelectorAll('.date-cell');

// Loop through each date cell and update its content
dateCells.forEach(cell => {
    const originalDate = cell.textContent;
    const formattedDate = formatDate(originalDate);
    cell.textContent = formattedDate;
});

//End change Date Format

// copy Code
function CopyToClipboard(text) {
    text=$(text).text();
    
    var $txt = $('<textarea />');

    $txt.val(text)
        .css({ width: "1px", height: "1px" })
        .appendTo('body');

    $txt.select();
    

    if (document.execCommand('copy')) {
        $txt.remove();
        alert("Copied");
    }
}
//Get currency sumbol from currency
function getCurrencySymbol(currencyCode) {
    alert(33333)
    // Create a new instance of Intl.NumberFormat
    const formatter = new Intl.NumberFormat('en', {
      style: 'currency',
      currency: currencyCode,
      minimumFractionDigits: 0, // Remove decimal places for the symbol
    });
  
    // Format a value and extract the currency symbol
    const parts = formatter.formatToParts(0);
    const symbol = parts.find(part => part.type === 'currency').value;

    alert(symbol)
  
    return symbol;
  }