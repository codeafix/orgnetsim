import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';
import { Results } from '../API/Results';
import { RunSpec } from '../API/Step';
import { ParseOptions } from '../API/Parse';
import { ResultsCsv } from '../API/Results';
import { SimList } from '../API/SimList';

const API = {
    rootPath: "http://localhost:8080",
    emptySimList: {simulations:[],notes:""},
    emptySim: function():SimInfo {return {id:"",name:"",description:"", steps:[], options:{
        linkTeamPeers: false,
        linkedTeamList: [],
        evangelistList: [],
        loneEvangelist: [],
        initColors: [],
        maxColors: 0,
        agentsWithMemory: false
    }};},
    sims: async function():Promise<SimList> {
        const response = await fetch(this.rootPath+"/api/simulation", {
        "method": "GET",
        "headers": {}
        }).catch(err => { 
            console.log(err);
        }) as Response;

        return response.json();
    },
    get: async function(id:string):Promise<SimInfo>{
        const response = await fetch(this.rootPath+"/api/simulation/"+id, {
            "method": "GET",
            "headers": {}
            }).catch(err => { console.log(err); 
        }) as Response;
        return response.json();
    },
    getStep: async function(path:string):Promise<Step>{
        const response = await fetch(this.rootPath+path, {
            "method": "GET",
            "headers": {}
            }).catch(err => { console.log(err); 
        }) as Response;
        return response.json();
    },
    getSteps: async function(sim:SimInfo):Promise<Array<Step>>{
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/step", {
            "method": "GET",
            "headers": {}
            }).catch(err => { console.log(err); 
        }) as Response;
        return response.json();
    },
    getResults: async function(sim:SimInfo):Promise<Results>{
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/results", {
            "method": "GET",
            "headers": {}
            }).catch(err => { console.log(err); 
        }) as Response;
        
        return response.json();
    },
    getResultsCsv: async function(sim:SimInfo):Promise<ResultsCsv>{
        var fn = "results.csv";
        const data = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/results", {
            "method": "GET",
            "headers": {
                "Content-Type": "text/csv"
            }
            }).then(response => {
                const contentDisp = response.headers.get('Content-Disposition') as string;
                const regExpFilename = /filename="(?<filename>[^"]*)"/;
                fn = regExpFilename.exec(contentDisp)?.groups?.filename ?? fn;
                return response.blob();
            }).then(blob => {
                return {
                    filename: fn,
                    blob: blob,
                };
            }).catch(err => { console.log(err); 
            }) as ResultsCsv;
        return data;
    },
    update: async function(sim:SimInfo):Promise<SimInfo>{
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id, {
            "method": "PUT",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(sim),
            })
            .catch(err => { console.log(err); 
            }) as Response;
        return response.json();
    },
    updateStep: async function(step:Step):Promise<Step>{
        const response = await fetch(this.rootPath+"/api/simulation/"+step.parent+"/step/"+step.id, {
            "method": "PUT",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(step),
            })
            .catch(err => { console.log(err); 
            }) as Response;
        return response.json();
    },
    runsim: async function(sim:SimInfo, spec:RunSpec){
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/run", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(spec),
            })
            .catch(err => { console.log(err); 
            }) as Response;
            return response.json();
    },
    parse: async function(sim:SimInfo, pdata:ParseOptions){
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/parse", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(pdata),
            }).catch(err => { console.log(err); 
            }) as Response;
            return response.json();
    },
    addlinks: async function(sim:SimInfo, pdata:ParseOptions){
        const response = await fetch(this.rootPath+"/api/simulation/"+sim.id+"/links", {
            "method": "PUT",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(pdata),
            }).catch(err => { console.log(err); 
            }) as Response;
            return response.json();
    },
    add: async function(name:string, description:string):Promise<SimInfo>{
        var sim = {name:name,description:description};
        const response = await fetch(this.rootPath+"/api/simulation", {
            "method": "POST",
            "headers": {
                'Content-Type': 'application/json'
            },
            "body": JSON.stringify(sim),
            })
            .catch(err => { console.log(err); 
            }) as Response;
        return response.json();
    },
    copy: async function(id:string):Promise<SimInfo>{
        const response = await fetch(this.rootPath+"/api/simulation/"+id+"/copy", {
            "method": "POST",
            "headers": {},
            })
            .catch(err => { console.log(err); 
            }) as Response;
        return response.json();
    },
    delete: async function(id:string){
        await fetch(this.rootPath+"/api/simulation/"+id, {
            "method": "DELETE",
            "headers": {}
            })
            .catch(err => { console.log(err); 
            }) as Response;
    },
}

export default API