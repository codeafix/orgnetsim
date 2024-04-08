import React, { useState, useEffect } from 'react';
import {Form, Modal, FormControl, Button, Table} from 'react-bootstrap'
import {X, Check} from 'react-bootstrap-icons';
import API from '../API/api';
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';
import { ParseOptions, Regex } from '../API/Parse';

type ParseSettingsProps = {
    sim:SimInfo;
    steps:Array<Step>;
    showimportmodal:boolean;
    readsim(id:string): void;
    onclose(): void;
}

const regexval = {
    // "2": "[1-9]\d*",

    // "3": "[1-9]\d*",

    "0":"^\\d+_(.*)$",
    "2":"^\\d+_(.*)$",
    "3":"\\S+",
}

const ParseSettings = (props:ParseSettingsProps) => {
    const [idcol, setidcol] = useState<number>(0);
    const [pcol, setpcol] = useState<number>(2);
    const [ncol, setncol] = useState<number>(0);
    const [delim, setdelim] = useState<string>(",");
    const [regex, setregex] = useState<Regex>(regexval);
    const [filetoupload, setfiletoupload] = useState<Blob>();
    const [showimpmodal, setshowimpmodal] = useState<boolean>(false);
    const [hasstep, sethasstep] = useState<boolean>(false);

    useEffect(() => {
        sethasstep((props.steps || []).length > 0)
        setshowimpmodal(props.showimportmodal);
    },[props.sim, props.steps, props.showimportmodal]);

    const handleimpclose = () => {
        setshowimpmodal(false);
        setidcol(0);
        setpcol(2);
        setncol(0);
        setdelim(",");
        setregex(regexval);
        setfiletoupload({} as Blob);
        props.onclose();
    };
    
    const handleimport = () => {
        const fr = new FileReader();
        if(!filetoupload?.text) return;

        fr.readAsDataURL(filetoupload);
        fr.onload = function() {
            if(!fr.result) return;
            const r = fr.result as string;
            const base64data = r.substring(r.indexOf(',') + 1);
        
            const pdata:ParseOptions = {
                "identifier": idcol,
                "parent": pcol,
                "name": ncol,
                "regex": regex,
                "delimiter": delim,
                "payload": base64data
            };
            if(hasstep){
                API.addlinks(props.sim, pdata).then(response => {
                    props.readsim(response.parent);
                });
            }else{
                API.parse(props.sim, pdata).then(response => {
                    props.readsim(response.parent);
                });
            }
        };
        handleimpclose();
    };
    
    return(
        <Modal
            show={showimpmodal}
            onHide={handleimpclose}
            backdrop="static"
            keyboard={false}
        >
            <Modal.Header closeButton>
                <Modal.Title>{hasstep?"Import Additional Links":"Import Network"}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form.Group controlId="form-file">
                    <Form.File label="Select File" onChange={(e:any) => {
                        if(e.target.files.length) setfiletoupload(e.target.files[0]);
                    }}/>
                </Form.Group>
                <Form.Group controlId="form-identifier">
                    <Form.Label>Identifier column</Form.Label>
                    <Form.Control type="number" value={idcol} onChange={e => setidcol((e.target as HTMLInputElement).valueAsNumber)}/>
                    <Form.Text className="text-muted">
                        The column that has the unique identifier in it
                    </Form.Text>
                </Form.Group>
                <Form.Group controlId="form-name" hidden={hasstep}>
                    <Form.Label>Name column</Form.Label>
                    <Form.Control type="number" value={ncol} onChange={e => setncol((e.target as HTMLInputElement).valueAsNumber)}/>
                    <Form.Text className="text-muted">
                        The column that has the agent name in it (this can be the same as the identifier column)
                    </Form.Text>
                </Form.Group>
                <Form.Group controlId="form-parent">
                    <Form.Label>Parent Identifier Column</Form.Label>
                    <Form.Control type="number" value={pcol} onChange={e => setpcol((e.target as HTMLInputElement).valueAsNumber)}/>
                    <Form.Text className="text-muted">
                        The column that has a parent identifier in it to show hierarchy in the network
                    </Form.Text>
                </Form.Group>
                <Form.Group className="mb-3" controlId="form-delimiter">
                    <Form.Label>Delimiter</Form.Label>
                    <FormControl type="string" value={delim} onChange={e => setdelim(e.target.value)}/>
                    <Form.Text className="text-muted">
                        The delimiter used to separate columns in rows of the imported file
                    </Form.Text>
                </Form.Group>
                <RegexControl regex={regex} onChange={(r:Regex) => setregex(r)}/>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="success" onClick={handleimport}>Save</Button>
                <Button variant="secondary" onClick={handleimpclose}>Cancel</Button>
            </Modal.Footer>
        </Modal>
    )
}

type RegexProps = {
    regex: Regex
    onChange: (regex: Regex) => void
}
type RegexRow = {
    col: string;
    regex: string;
}

const RegexControl = (props:RegexProps) => {
    const [r, setr] = useState<Array<RegexRow>>([]);
    const [showaddrow, setShowaddrow] = useState<boolean>(false);
    const [newcol, setnewcol] = useState<string>("");
    const [newregex, setnewregex] = useState<string>("");

    useEffect(() => {
        var t = Object.entries(props.regex).map(([k,v]) => {return {col:k, regex:v}});
        setr(t);
    },[props.regex]);

    const addRegex = () => {
        var ir = r.findIndex(row => row.col === newcol);
        var t = r.slice();
        if(ir >= 0){
            t[ir].regex = newregex;
        }else{
            t.push({col:newcol, regex:newregex});
        }
        setnewcol("");
        setnewregex("");
        raiseChange(t);
        setShowaddrow(false);
    }

    const removeRegex = (i:number) => {
        var t = r.toSpliced(i,1);
        raiseChange(t);
    }

    const raiseChange = (t:Array<RegexRow>) => {
        props.onChange(t.reduce((acc, row) => {acc[row.col] = row.regex; return acc}, {} as Regex));
    }

    return (
        <Form.Group className="mb-3" controlId="form-regex">
            <Form.Label>Column Regular Expression</Form.Label>
            <Table bordered size="sm">
                <thead>
                    <tr>
                        <th key="h_iter">Column</th>
                        <th key="h_conv">Regex</th>
                    </tr>
                </thead>
                <tbody>
                    {r.map((regexrow, i) => {
                        return <tr key={i}>
                            <td key={"col_"+i}>{regexrow.col}</td>
                            <td key={"regex_"+i}>{regexrow.regex}<Button className="p-0 border-0 btn btn-light float-right" onClick={() => {removeRegex(i)}}><X/></Button></td>
                        </tr>
                    })} 
                    <tr key="add_row" hidden={!showaddrow}>
                        <td><FormControl type="string" value={newcol} onChange={e => setnewcol(e.target.value)}/></td>
                        <td>
                            <FormControl className="float-left w-auto" type="string" value={newregex} onChange={e => setnewregex(e.target.value)}/>
                            <Button className="p-0 border-0 btn btn-light float-right" onClick={() => setShowaddrow(false)}><X/></Button>
                            <Button className="p-0 border-0 btn btn-light float-right" onClick={() => {addRegex()}}><Check/></Button></td>
                    </tr>
                </tbody>
            </Table>
            <Button size="sm" className="btn btn-primary float-right" onClick={() => setShowaddrow(true)}>Add</Button>
            <Form.Text className="text-muted">
                Regular expressions applied to each column to decide whether to include a row in the import
            </Form.Text>
        </Form.Group>
    )
}

export default ParseSettings