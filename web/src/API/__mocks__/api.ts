import { SimInfo } from '../../API/SimInfo';
import { Step } from '../../API/Step';
import { Results } from '../../API/Results';
import { RunSpec } from '../../API/Step';
import { ParseOptions } from '../../API/Parse';
import { ResultsCsv } from '../../API/Results';
import { SimList } from '../../API/SimList';
import { Network } from '../../API/Network';
import { AgentState } from '../../API/Network';

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
    getStepNetwork: async function(simId: string, stepId: string): Promise<Network> {
        return {"links":[{"source":"Lead","target":"Dev_3","strength":35},{"source":"Lead","target":"Dev_4","strength":35},{"source":"Lead","target":"Dev_5","strength":40},{"source":"Lead","target":"Dev_1","strength":45},{"source":"Lead","target":"Dev_2","strength":41},{"source":"Dev_3","target":"Dev_4","strength":47},{"source":"Dev_3","target":"Dev_5","strength":45},{"source":"Dev_3","target":"Dev_1","strength":36},{"source":"Dev_3","target":"Dev_2","strength":37},{"source":"Dev_4","target":"Dev_5","strength":34},{"source":"Dev_4","target":"Dev_1","strength":39},{"source":"Dev_4","target":"Dev_2","strength":37},{"source":"Dev_5","target":"Dev_1","strength":34},{"source":"Dev_5","target":"Dev_2","strength":45},{"source":"Dev_1","target":"Dev_2","strength":38}],"nodes":[{"id":"Lead","name":"Lead","color":1,"susceptability":0.5737704957356361,"influence":1.0983510095299653,"contrariness":0.7434660073333511,"change":45,"type":"AgentState","fx":512.0238945054554,"fy":383.2777626590721},{"id":"Dev_1","name":"Dev 1","color":1,"susceptability":0.8666370569034374,"influence":0.7173451361458638,"contrariness":0.44984596307027497,"change":13,"type":"AgentState","fx":498.7319761723783,"fy":406.0009971245655},{"id":"Dev_2","name":"Dev 2","color":0,"susceptability":1.177164998572346,"influence":1.0537797825980781,"contrariness":0.940969259983882,"change":0,"type":"AgentState","fx":509.9407199938294,"fy":358.1175533021693},{"id":"Dev_3","name":"Dev 3","color":0,"susceptability":1.123690974981149,"influence":0.5847478097190408,"contrariness":0.8289627855715872,"change":0,"type":"AgentState","fx":528.3743371463426,"fy":403.6228601629286},{"id":"Dev_4","name":"Dev 4","color":0,"susceptability":0.8812168832676177,"influence":0.7392639453582498,"contrariness":0.4325477027577978,"change":18,"type":"AgentState","fx":486.83099078103857,"fy":378.383316333952},{"id":"Dev_5","name":"Dev 5","color":0,"susceptability":1.2810755968270868,"influence":1.0332774519895156,"contrariness":0.7011268127170646,"change":0,"type":"AgentState","fx":536.0980539104477,"fy":374.5974462319489}],"maxColors":2};
    },
    getStepAgents: async function(simId: string, stepId: string): Promise<AgentState[]> {
        return [{"id":"Lead","name":"Lead","color":1,"susceptability":0.5737704957356361,"influence":1.0983510095299653,"contrariness":0.7434660073333511,"change":45,"type":"AgentState","fx":512.0238945054554,"fy":383.2777626590721},{"id":"Dev_1","name":"Dev 1","color":1,"susceptability":0.8666370569034374,"influence":0.7173451361458638,"contrariness":0.44984596307027497,"change":13,"type":"AgentState","fx":498.7319761723783,"fy":406.0009971245655},{"id":"Dev_2","name":"Dev 2","color":0,"susceptability":1.177164998572346,"influence":1.0537797825980781,"contrariness":0.940969259983882,"change":0,"type":"AgentState","fx":509.9407199938294,"fy":358.1175533021693},{"id":"Dev_3","name":"Dev 3","color":0,"susceptability":1.123690974981149,"influence":0.5847478097190408,"contrariness":0.8289627855715872,"change":0,"type":"AgentState","fx":528.3743371463426,"fy":403.6228601629286},{"id":"Dev_4","name":"Dev 4","color":0,"susceptability":0.8812168832676177,"influence":0.7392639453582498,"contrariness":0.4325477027577978,"change":18,"type":"AgentState","fx":486.83099078103857,"fy":378.383316333952},{"id":"Dev_5","name":"Dev 5","color":0,"susceptability":1.2810755968270868,"influence":1.0332774519895156,"contrariness":0.7011268127170646,"change":0,"type":"AgentState","fx":536.0980539104477,"fy":374.5974462319489}];
    },
    updateStepNetwork: async function(simId: string, stepId: string, network: Network): Promise<Network> {
        return network;
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