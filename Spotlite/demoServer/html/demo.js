// Markers and other controls are stored in the variables

// variable for map
var map;
// marker, info window variables

var markers = [];
var overlapmarkers = [];

var rider_icon = 'images/taxi_call.png'
var cab_icon = 'images/black_taxi.png'
var res_cab_icon = 'images/orange_car.gif'

// Initialize the main map
function initializeMap() {
    console.log("Initializing the map")

    // Create the SF map with the center of city as initial value
    var sfcity = new google.maps.LatLng(37.78, -122.454150)

    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 13,
        center: sfcity,
        mapTypeId: google.maps.MapTypeId.ROADMAP
    });

    curRadius = 0.6215;

    // Anchor marker ... user's selection, or typed in the address bar
    // Create the search box and link it to the UI element.
    var defaultBounds = new google.maps.LatLngBounds(
        new google.maps.LatLng(37.5902, -122.1759),
        new google.maps.LatLng(37.8474, -122.5631)
    );


    map.addListener('click', function(e) {
        checkOverlap(e.latLng)
    });

    map.addListener('rightclick', function(e) {
        console.log(e.latLng.lat(), e.latLng.lng())
    });


    drawForDemo()
        //drawTestLog()

    //setInterval(ajaxCall, 6000); //300000 MS == 5 minutes


}
/*
function ajaxCall() {
    //do your AJAX stuff here
    var searchurl = "http://ec2-50-112-18-158.us-west-2.compute.amazonaws.com:5000/livemap"
    $.ajax({
        type: 'GET',
        url: searchurl,
        data: {},
        dataType: 'json',
        success: function(data) {
            console.log(data)
        },
        error: function() {}
    });

}
*/

function checkOverlap(latLng) {

    var searchurl = "/Overlap?lat=" + latLng.lat() + "&lng=" + latLng.lng() + "&";
    $.ajax({
        type: 'GET',
        url: searchurl,
        data: {},
        dataType: 'json',
        success: function(data) {
            drawOverlapPoints(latLng.lat(), latLng.lng(), data)
        },
        error: function() {}
    });
}

/*
function doSearchAndUpdate(latLng) {
    circle_marker.setCenter(latLng)
    deleteMarkers()
    var searchurl = "http://127.0.0.1:4747/CircleCells?lat=" + latLng.lat() + "&lng=" + latLng.lng() + "&";
    $.ajax({
        type: 'GET',
        url: searchurl,
        data: {},
        dataType: 'json',
        success: function(data) {
            drawCells(data)
        },
        error: function() {
            alert('Error occured while getting food truck locations !!!');
        }
    });

    var searchurl2 = "http://127.0.0.1:4747/PointCells?lat=" + latLng.lat() + "&lng=" + latLng.lng() + "&";
    $.ajax({
        type: 'GET',
        url: searchurl2,
        data: {},
        dataType: 'json',
        success: function(data) {
            drawAnchorCell(data)
        },
        error: function() {
            alert('Error occured while getting food truck locations !!!');
        }
    });

}
*/

function drawRider(lat, lng) {
    var latlng = new google.maps.LatLng(lat, lng)

    var ridermarker = new google.maps.Marker({
        position: latlng,
        map: map,
        icon: rider_icon,
        clickable: false
    });

    //ridermarker.setMap(map);
    overlapmarkers.push(ridermarker);
}

function drawForDemo() {
    drawDriverCircle(233, 37.775760, -122.412360, 2.346178)
    drawDriverCircle(234, 37.790050, -122.403840, 0.000000)
    drawDriverCircle(235, 37.687110, -122.451570, 0.953019)
    drawDriverCircle(236, 37.677030, -122.388410, 2.257210)
    drawDriverCircle(237, 37.783710, -122.424370, 1.000000)
    drawDriverCircle(238, 37.768670, -122.427360, 0.000000)
    drawDriverCircle(239, 37.749700, -122.396470, 3.856807)
    drawDriverCircle(240, 37.775210, -122.424350, 0.000000)
}

function drawDriverCircle(id, lat, lng, radius) {
    var dloc = new google.maps.LatLng(lat, lng)
    drawDriver(dloc, radius)
}

/*

function drawClosestDriver(slat, slng, elat, elng) {
    var s = new google.maps.LatLng(slat, slng)
    var e = new google.maps.LatLng(elat, elng)
    var coords = [
        s, e
    ];

    var linemarker = new google.maps.Polyline({
        path: coords,
        strokeColor: '#424242',
        strokeOpacity: 1.0,
        strokeWeight: 2
    });

    linemarker.setMap(map);
    markers.push(linemarker);
}
*/

function drawDriver(latLng, radius) {

    if (radius === 0.0) {
        var cabmarker = new google.maps.Marker({
            position: latLng,
            map: map,
            icon: cab_icon,
            clickable: false
        });
        markers.push(cabmarker);

    } else {
        var cabmarker = new google.maps.Marker({
            position: latLng,
            map: map,
            icon: res_cab_icon,
            clickable: false
        });
        markers.push(cabmarker)
    }


    var cmarker = new google.maps.Circle({
        strokeColor: '#B18904',
        strokeOpacity: 0.95,
        strokeWeight: 2,
        fillColor: '#F7D358',
        fillOpacity: 0.45,
        map: map,
        radius: radius * 1000.0,
        clickable: false
    });

    cmarker.setMap(map);
    cmarker.setCenter(latLng);
    markers.push(cmarker);


}


function drawArrow(obj, platLng) {
    /*
    console.log(obj)
    var cmarker = new google.maps.Circle({
        strokeColor: '#8A0829',
        strokeOpacity: 0.95,
        strokeWeight: 1,
        fillColor: '#8A0829',
        fillOpacity: 0.70,
        map: map,
        zIndex: 10,
        radius: 75.0,
        clickable: false
    });



    cmarker.setMap(map);
    cmarker.setCenter(latLng);
    overlapmarkers.push(cmarker)
*/
    // draw a line ...

    var latLng = new google.maps.LatLng(obj.Lat, obj.Lng)
    var coords = [
        platLng, latLng
    ];

    var routermarker = new google.maps.Polyline({
        path: coords,
        strokeColor: '#424242',
        strokeOpacity: 1.0,
        strokeWeight: 2
    });

    routermarker.setMap(map);
    overlapmarkers.push(routermarker)

}

/*
function drawS2Cell(obj, color) {
    var coords = []
    for (var i = 0; i < obj.length; i++) {
        var loc = {
            lat: obj[i].Lat,
            lng: obj[i].Lng
        };
        coords.push(loc)
    }

    // Construct the polygon.
    var s2marker = new google.maps.Polygon({
        paths: coords,
        strokeColor: color,
        strokeOpacity: 0.8,
        strokeWeight: 1,
        fillColor: color,
        fillOpacity: 0.40
    });
    s2marker.setMap(map);

    markers.push(s2marker);
    return s2marker;
}
*/


// Sets the map on all markers in the array.
function setMapOnAll(map) {
    for (var i = 0; i < markers.length; i++) {
        markers[i].setMap(map);
    }
}

// Removes the markers from the map, but keeps them in the array.
function clearMarkers() {
    setMapOnAll(null);
}

// Shows any markers currently in the array.
function showMarkers() {
    setMapOnAll(map);
}

function deleteMarkers() {
    clearMarkers();
    markers = [];
}


function drawCells(data, color) {

    if (data) {
        for (var i = 0; i < data.length; i++) {
            var obj = data[i];
            drawCell(obj, color)
        }
    }
}

function drawOverlapPoints(lat, lng, data) {

    for (var i = 0; i < overlapmarkers.length; i++) {
        overlapmarkers[i].setMap(null);
    }
    overlapmarkers = []


    drawRider(lat, lng)

    var latLng = new google.maps.LatLng(lat, lng)
        /*
                    var personmarker = new google.maps.Marker({
                        position: latLng,
                        map: map,
                        icon: rider_icon,
                        clickable: false
                    });
                    overlapmarkers.push(personmarker)
                */

    if (data) {
        for (var i = 0; i < data.length; i++) {
            var obj = data[i];
            drawArrow(obj, latLng)
        }
    }
}