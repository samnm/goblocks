#version 120

uniform float timer;
uniform mat4 mv_matrix;
uniform mat4 p_matrix;

attribute vec3 position;

varying vec2 texcoord;
varying float fade_factor;

void main()
{
    vec4 eye_position = mv_matrix * vec4(position, 1.0);
    gl_Position = p_matrix * eye_position;
    texcoord = position.xy + vec2(0.5);
    fade_factor = sin(timer) * 0.5 + 0.5;
}
