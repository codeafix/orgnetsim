import React, { useEffect, useRef } from 'react';
import { scaleLinear, max, select } from 'd3';
import Color from './Color';
import API from '../api';

const AgentColorChart = (props) => {
    const chart = useRef(null);

    const createChart = (results) => {
        const width = chart.current.parentElement.offsetWidth,
            height = Math.round(width/1.6),
            rMargin = 40, bMargin = 30;
        
        const resize = () => {
            const w = chart.current.parentElement.offsetWidth - rMargin;
            select(chart.current).attr('width', w)
                .attr('height', Math.round(w / 1.6));
        };

        select(window).on(
            'resize.' + select(chart.current.parentElement).attr('id'), 
            resize
        );
        
        select(chart.current).attr('viewBox', `0 0 ${width} ${height}`)
            .attr('preserveAspectRatio', 'xMinYMid')
            .call(resize);

        const dataMax = max(results['conversations']);
        const yScale = scaleLinear()
            .domain([0, dataMax])
            .range([0, height-bMargin]);

        select(chart.current)
            .selectAll('rect')
            .data(results['conversations'])
            .enter()
            .append('rect')
         
        select(chart.current)
            .selectAll('rect')
            .data(results['conversations'])
            .exit()
            .remove()
        
        select(chart.current)
            .selectAll('rect')
            .data(results['conversations'])
            .style('fill', Color.cssColorFromVal(0))
            .attr('x', (d,i) => i)
            .attr('y', d => height - bMargin - yScale(d))
            .attr('height', d => yScale(d))
            .attr('width', 1)
    }

    useEffect(() => {
        if (!props.sim['id']) {
            return
        }
        API.getResults(props.sim).then(results => {
            createChart(results);
        }).catch(err => console.error(err));
    },[props.sim]);

    return <svg ref={chart}/>
}

export default AgentColorChart;