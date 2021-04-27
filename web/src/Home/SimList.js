import React from 'react';
import {Link} from 'react-router-dom';

const SimList = ({sims}) => {
    return sims.map(sim => {
        return(
            <div class="card" key={sim.id}>
                <Link class = "card-header" to={'/simulation/' + sim.id}>{sim.name}</Link>
                <div class = "card-body">{sim.description}</div>
            </div>
        );
    })
}

export default SimList