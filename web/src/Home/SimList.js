import React from 'react';
import {Link} from 'react-router-dom';

const SimList = ({sims}) => {
    return sims.map(sim => {
        return(
            <div class="card" key={sim.ID}>
                <Link class = "card-header" to={'/simulation/' + sim.ID}>{sim.Name}</Link>
                <div class = "card-body">{sim.Description}</div>
            </div>
        );
    })
}

export default SimList