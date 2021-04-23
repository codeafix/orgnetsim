import React, { Component } from 'react';
import API from '../api';
import SimList from './SimList';

class Home extends Component {
    constructor(props){
        super(props);
        this.state = {
            simulations: [],
        };
    }

    componentDidMount(){
        this.setState({
            simulations: API.sims(),
        });
    }

    addSimulation(){
        const sim = API.add();
        const sims = this.state.simulations.concat(sim);
        this.setState({
            simulations: sims,
        });
    }

    render(){
        const sims = this.state.simulations
        return(
            <div>
                <h1>Simulation Set</h1>
                <p>{API.notes}</p>
                <h2>List of Simulations</h2>
                <SimList sims={sims}/>
                <button onClick={() => this.addSimulation()}>Add</button>
            </div>
        )
    }
}

export default Home