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

        // build image chart
        $.get(`/api/v1/timeline?img=${img}`, function(d) {
            console.log(d);
            var imgChart = c3.generate({
                bindto: '#image-chart',
                data: {
                    json: d.data,
                    keys: {
                        x: 'date',
                        value: ['grype', 'trivy', 'snyk'],
                    },
                    type: 'line',
                    labels: true
                },
                size: {
                    height: 200
                },
                padding: {
                    right: 20,
                    bottom: 0
                },
                axis: {
                  x: {
                    x: ['date'],
                    type: "timeseries",
                    label: {
                        show: false
                    }
                  },
                  y: {
                    label: {
                        show: false
                    },
                    tick: {
                        format: function (d) {
                            return (parseInt(d) == d) ? d : null;
                        },
                    }
                  }
                },
            });
        });

        // build vulnerabilities chart
        var vulChart = c3.generate({
            bindto: '#vuln-chart',
            data: {
                columns: vulnData,
                type : 'pie',
            },
            size: {
                width: 300
            },
            legend: {
                position: 'right'
            }
        });
    }

});