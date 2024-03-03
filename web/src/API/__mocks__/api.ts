import { SimInfo } from '../../API/SimInfo';
import { Step } from '../../API/Step';
import { Results } from '../../API/Results';
import { RunSpec } from '../../API/Step';
import { ParseOptions } from '../../API/Parse';
import { ResultsCsv } from '../../API/Results';
import { SimList } from '../../API/SimList';

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
        // return this.emptySimList;
        return {notes:"Some notes",simulations:[
            {id:"1",name:"test1",description:"test1desc",steps:[],options:{
                linkTeamPeers: false,
                linkedTeamList: [],
                evangelistList: [],
                loneEvangelist: [],
                initColors: [],
                maxColors: 0,
                agentsWithMemory: false
            }},
            {id:"2",name:"test2",description:"test2desc",steps:[],options:{
                linkTeamPeers: false,
                linkedTeamList: [],
                evangelistList: [],
                loneEvangelist: [],
                initColors: [],
                maxColors: 0,
                agentsWithMemory: false
            }}
        ]};
    },
    get: async function(id:string):Promise<SimInfo>{
        return this.emptySim();
    },
    getStep: async function(path:string):Promise<Step>{
        return {
            id: "1",
            parent: "1",
            network:{
                nodes: [],
                links: [],
                maxColors: 0,
            },
            results:{
                iterations: 0,
                colors: [[]],
                conversations: []
            }
        };
    },
    getSteps: async function(sim:SimInfo):Promise<Array<Step>>{
        return [];
    },
    getResults: async function(sim:SimInfo):Promise<Results>{
        return {"iterations":100,"colors":[[6,0],[4,2],[6,0],[5,1],[4,2],[4,2],[4,2],[5,1],[5,1],[5,1],[5,1],[6,0],[4,2],[4,2],[4,2],[4,2],[5,1],[4,2],[4,2],[5,1],[5,1],[6,0],[5,1],[6,0],[6,0],[6,0],[6,0],[6,0],[5,1],[5,1],[5,1],[6,0],[4,2],[5,1],[4,2],[5,1],[5,1],[5,1],[6,0],[5,1],[4,2],[4,2],[6,0],[5,1],[5,1],[5,1],[6,0],[5,1],[5,1],[5,1],[5,1],[4,2],[5,1],[5,1],[5,1],[4,2],[5,1],[6,0],[6,0],[6,0],[6,0],[5,1],[5,1],[5,1],[5,1],[6,0],[4,2],[5,1],[6,0],[4,2],[5,1],[5,1],[5,1],[5,1],[4,2],[5,1],[4,2],[5,1],[4,2],[5,1],[6,0],[5,1],[5,1],[6,0],[6,0],[5,1],[5,1],[6,0],[4,2],[4,2],[5,1],[3,3],[3,3],[4,2],[3,3],[4,2],[4,2],[4,2],[6,0],[4,2],[4,2]],"conversations":[0,6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,5,6,5,6,6,6,6,5,6,6,6,6,6,6]};
    },
    getResultsCsv: async function(sim:SimInfo):Promise<ResultsCsv>{
        return {filename:"results.csv",blob:new Blob()};
    },
    update: async function(sim:SimInfo):Promise<SimInfo>{
        return sim;
    },
    updateStep: async function(step:Step):Promise<Step>{
        return step;
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
    add: async function(name:string, description:string):Promise<SimInfo>{
        return {id:"1",name:name,description:description, steps:[], options:{
            linkTeamPeers: false,
            linkedTeamList: [],
            evangelistList: [],
            loneEvangelist: [],
            initColors: [],
            maxColors: 0,
            agentsWithMemory: false
        }};
    },
    delete: async function(id:string){
        return;
    },
}

export default API