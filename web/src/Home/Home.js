import React, { Component } from 'react';
import API from '../api';
import SimList from './SimList';

class Home extends Component {
    constructor(props){
        super(props);
        this.state = {
            simlist: API.emptySimList,
        };
    }

    componentDidMount(){
        API.sims()
            .then(response => {
                this.setState({
                    simlist: response,
                });
                API.simCount = response.simulations.length
            })
    }

    addSimulation(){
        API.add().then(response => {
            const simlist = this.state.simlist;
            simlist.simulations = simlist.simulations.concat(response);
            this.setState({
                simlist: simlist,
            });
        });
    }

    render(){
        const sims = this.state.simlist.simulations
        const notes = this.state.simlist.notes
        return(
            <div>
                <h1>Simulation Set</h1>
                <p>{notes}</p>
                <h2>List of Simulations</h2>
                <SimList sims={sims}/>
                <button onClick={() => this.addSimulation()}>Add</button>
            </div>
        )
    }
}

export default Home