const urlParams = new URLSearchParams(window.location.search);

$(function () {

    // === HOME PAGE ===

    // Autocomplete for image search
    if($( "#image-search" ).length){
        $( "#image-search" ).autocomplete({  
            source: images,
            autoFocus: true,
            select: function( event, ui ) {
                console.log(ui.item.value);
                uri = ui.item.value;
                window.location.href = '/image?img=' + encodeURIComponent(uri);
            }
        }); 
    }

    // === IMAGE PAGE ===
    if($( "#image-chart" ).length){
        var img = urlParams.get('img');
        $.get(`/api/v1/timeline?img=${img}`, function(d) {
            console.log(d);
            var imgChart = c3.generate({
                bindto: '#image-chart',
                data: {
                    json: d.data,
                    keys: {
                        x: 'date',
                        value: ['total'],
                    },
                    type: 'line',
                    labels: true
                },
                axis: {
                  x: {
                    x: ['source'],
                    type: "timeseries",
                    label: 'Vulnerabilities'
                  },
                  y2: {
                    type: 'timeseries',
                    label: {
                        text: 'Source',
                    }
                  }
                },
            });
        });
    }

});