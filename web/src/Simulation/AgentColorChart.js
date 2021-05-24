import React, { useEffect, useRef, useState } from 'react';
import { scaleLinear, max, select, stack, axisBottom, axisLeft, area } from 'd3';
import Color from './Color';
import API from '../api';
import Spinner from 'react-bootstrap/Spinner';

const AgentColorChart = (props) => {
    const chart = useRef(null);
    const [loading, setloading] = useState(false);

    const createChart = (results) => {
        if(!results['colors'][0]) return;
        const maxColors = props.sim.options['maxColors'];
        const margin = {top: 10, right: 60, bottom: 20, left: 40},
            vwidth = chart.current.parentElement.offsetWidth,
            cw = vwidth - margin.left - margin.right,
            vheight = Math.round(cw/1.6),
            ch = vheight - margin.top - margin.bottom;
        
        const resize = () => {
            const w = chart.current.parentElement.offsetWidth,
                h = Math.round((w -  margin.left - margin.right)/1.6);
            select(chart.current).attr('width', w)
                .attr('height', h);
        };

        select(window).on(
            'resize.' + select(chart.current.parentElement).attr('id'), 
            resize
        );
        
        const svg = select(chart.current);

        svg.attr('viewBox', `0 0 ${vwidth} ${vheight}`)
            .attr('preserveAspectRatio', 'xMinYMid')
            .call(resize);

        const dataMax = results['colors'][0].reduce((a, b) => a + b, 0);
        const iterations = results['iterations'];
        const chartData = results['colors'];

        var stackedData = stack()
            .keys(Color.colorValSlice(maxColors))
            .value(function(d, key){
                return d[key]
            })
            (chartData)
            
        //X Axis
        var xh = ch + margin.top;
        var xScale = scaleLinear()
            .domain([0, iterations])
            .range([0, cw]);
        svg.append("g")
            .attr("class", "small")
            .attr("transform", "translate(" + margin.left + "," + xh + ")")
            .call(axisBottom(xScale).ticks(10));
        
        //Add Y axis
        var yScale = scaleLinear()
            .domain([0, dataMax])
            .range([ch, 0]);
        svg.append("g")
            .attr("class", "small")
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
            .call(axisLeft(yScale).ticks(10));
          
        svg.selectAll("mylayers")
            .data(stackedData)
            .enter()
            .append("path")
            .style("fill", (d) => Color.cssColorFromVal(d.key))
            .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
            .attr("d", area()
                .x(function(d, i) { return xScale(i); })
                .y0(function(d) { return yScale(d[0]); })
                .y1(function(d) { return yScale(d[1]); })
            );
    }

    useEffect(() => {
        if (!props.sim['id']) {
            return
        }
        setloading(true);
        select(chart.current).selectAll("*").remove();
        API.getResults(props.sim).then(results => {
            createChart(results);
            setloading(false);
        }).catch(err => {
            console.error(err);
            setloading(false);
        });
    },[props.sim]);

    return <div>
            {loading && <Spinner animation="border" variant="info" />}
            <svg class="mb-3" ref={chart}/>
        </div>
}

export default AgentColorChart;