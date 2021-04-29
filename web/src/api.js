const API = {
    rootPath: "http://localhost:8080/api/",
    emptySimList: {simulations:[],notes:""},
    sims: async function() {
        const response = await fetch(this.rootPath+"simulation", {
        "method": "GET",
        "headers": {}
        })
        .catch(err => { console.log(err); 
        });
        return response.json();
    },
    get: async function(id){
        const response = await fetch(this.rootPath+"simulation/"+id, {
            "method": "GET",
            "headers": {}
            })
            .catch(err => { console.log(err); 
            });
            return response.json();
    },
    add: async function(name, description){
        var sim = {name:name,description:description};
        const response = await fetch(this.rootPath+"simulation", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(sim),
            })
            .catch(err => { console.log(err); 
            });
        return response.json();
    },
}

export default API