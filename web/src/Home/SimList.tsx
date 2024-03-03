import React from 'react';
import {Link} from 'react-router-dom';
import {Card} from 'react-bootstrap';
import {Trash} from 'react-bootstrap-icons';
import { SimInfo } from '../API/SimInfo';


type SimListProps = {
    deleteFunc(id: string, name: string): void;
    sims: Array<SimInfo>;
}

const SimList = (props:SimListProps) => {

    return props.sims.map(sim => {
            return(
                <Card border="info" style={{ width: '18rem' }} key={sim.id}>
                    <Card.Header><Card.Title>{sim.name}<Trash className="float-right" onClick={() => props.deleteFunc(sim.id, sim.name)}/></Card.Title></Card.Header>
                    <Card.Body>
                        <Card.Text>{sim.description}</Card.Text>
                    </Card.Body>
                    <Card.Footer>
                        <Link className="btn btn-primary" role="button" to={'/simulation/' + sim.id}>Open</Link>
                    </Card.Footer>
                </Card>
            );
        });
}

export default SimList