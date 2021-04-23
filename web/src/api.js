const API = {
    simArray: [
        {ID: 1, Name: "sim 1", Description: "The first sim"},
        {ID: 2, Name: "sim 2", Description: "The second sim"},
        {ID: 3, Name: "sim 3", Description: "The third sim"},
        {ID: 4, Name: "sim 4", Description: "The fourth sim"},
        {ID: 5, Name: "sim 5", Description: "The fifth sim"},
    ],
    notes: "Here are the notes describing what this simulation set is studying.",
    sims: function() {return this.simArray.slice()},
    get: function(id){
        const isSim = p => p.ID === id
        return this.simArray.find(isSim)
    },
    add: function(){
        var n = this.simArray.length+1
        var sim = {ID:n, Name:"sim "+n,Description:"Simulation number "+n};
        this.simArray.push(sim);
        return sim;
    },
}

export default API