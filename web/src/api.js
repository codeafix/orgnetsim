const API = {
    rootPath: "http://localhost:8080",
    emptySimList: {simulations:[],notes:""},
    sims: async function() {
        const response = await fetch(this.rootPath+"/api/simulation", {
        "method": "GET",
        "headers": {}
        })
        .catch(err => { console.log(err); 
        });
        return response.json();
    },
    get: async function(id){
        const response = await fetch(this.rootPath+"/api/simulation/"+id, {
            "method": "GET",
            "headers": {}
            })
            .catch(err => { console.log(err); 
            });
            return response.json();
    },
    getStep: async function(path){
        const response = await fetch(this.rootPath+path, {
            "method": "GET",
            "headers": {}
            })
            .catch(err => { console.log(err); 
            });
            return response.json();
    },
    update: async function(sim){
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id, {
            "method": "PUT",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(sim),
            })
            .catch(err => { console.log(err); 
            });
            return response.json();
    },
    runsim: async function(sim, spec){
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/run", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(spec),
            })
            .catch(err => { console.log(err); 
            });
            return response.json();
    },
    parse: async function(sim, pdata){
        const response = await fetch("http://localhost:8080/api/simulation/"+sim.id+"/parse", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(pdata),
            }).catch(err => { console.log(err); 
            });
            return response.json();
    },
    add: async function(name, description){
        var sim = {name:name,description:description};
        const response = await fetch(this.rootPath+"/api/simulation", {
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
    delete: async function(id){
        await fetch(this.rootPath+"/api/simulation/"+id, {
            "method": "DELETE",
            "headers": {}
            })
            .catch(err => { console.log(err); 
            });
    },
}

export default API