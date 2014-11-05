cmddata.push({
    "order_id": k,
    "count_cost": sp.toString()
});


var cmd = {
    "oper_type": "3",
    "fmt_ver": "1",
    "data": cmddata
};

readStream.on('data', function(data) {

    var res = String.fromCharCode(02);

    console.log(data.toString());

    data = data.toString().split(res);

    if(data.length < 3 ){
        return ;
    }

    var price = 0;

    var adid = data[3];
    var orderid = data[4];

    if (data[1] == "5" || data[1] == "6" || data[1] == "24") {

        price = data[12];

    } else if (data[1] == "8") {
        price = data[20];

    } else if (data[1] == "16", data[1] == "19", data[1] == "27") {

        price = data[20];
    }

    pricelog.push({
        orderid: orderid,
        price: parseInt(price)
    });

    //console.log(pricelog);

});
