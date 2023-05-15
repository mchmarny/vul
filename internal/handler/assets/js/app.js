const urlParams = new URLSearchParams(window.location.search);
const cssHighlight = 'highlight';

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
                color: {
                    pattern: ['#0a036b', '#fbb504', '#7b0265']
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
            color: {
                /*
                Negligible, Low, Medium, High, Critical, Unknown
                */
                pattern: ['#FFDCB6', '#F79540', '#FC4F00', '#B71375', '#8B1874', '#9BA4B5']
            },
            size: {
                width: 400
            },
            legend: {
                position: 'right'
            }
        });
    }

    // === VULNERABILITIES PAGE ===
    if($( '.exposure-nav' ).length){

        var countFunc = function() {
            $('#pkg-counter').remove();
            var count = $( '.package:visible' ).length;
            jQuery('<div>', {class: 'ui-helper-hidden-accessible', id: 'pkg-counter'})
                .html(count + ' packages found')
                .appendTo('body');

            setTimeout(function(){
                $('#pkg-counter').remove();
            }, 5000);
        }

        // unique source toggle
        $( '#unique' ).click(function(e) {
            e.preventDefault();
            $( '.same' ).parent().parent().hide();
            countFunc();
        });

        // fixed version taggle
        $( '#fixed' ).click(function(e) {
            e.preventDefault();
            $( '.nofixed' ).parent().parent().parent().hide();
            countFunc();
        });

        // show all 
        $( '#show-all' ).click(function(e) {
            e.preventDefault();
            $( '.same' ).parent().parent().show();
            $( '.nofixed' ).parent().parent().parent().show();
            countFunc();
        });

        // filter by source
        $( "#vul-filter" ).on( "keyup" , function() {
            var value = $(this).val().toLowerCase();
            if (value === '') {
                $( '.package-title' ).removeClass( cssHighlight );
                $( '.exposure a' ).parent().removeClass( cssHighlight );
                return;
            }
        
            $( '.package-title' ).each(function(index) {
                $(this).removeClass( cssHighlight );
                var text = $(this).text().toLowerCase();
                if (text.indexOf(value) !== -1) {
                    $(this).addClass( cssHighlight );
                }
            });

            $( '.exposure a' ).each(function(index) {
                $(this).parent().removeClass( cssHighlight );
                var text = $(this).text().toLowerCase();
                if (text.indexOf(value) !== -1) {
                    $(this).parent().addClass( cssHighlight );
                }
            });
        });
    }

});