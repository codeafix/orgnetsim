import React, { useState, useEffect } from 'react';
import API from '../api';
import SimList from './SimList';

const Home = () => {
    const [simlist, setSimlist] = useState([]);
    const [notes, setNotes] = useState("");
    
    useEffect(() => {
        API.sims()
            .then(response => {
                setSimlist(response.simulations);
                setNotes(response.notes)
                API.simCount = response.simulations.length;
            })
      },[]);

    function addSimulation(){
        API.add().then(response => {
            setSimlist(simlist.concat(response));
        });
    }
    
    return(
        <div>
            <h1>Simulation Set</h1>
            <p>{notes}</p>
            <h2>List of Simulations</h2>
            <SimList sims={simlist}/>
            <button onClick={addSimulation}>Add</button>
        </div>
    )
}

export default Home