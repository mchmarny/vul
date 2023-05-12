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

});