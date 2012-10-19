#version 120

uniform float timer;
uniform mat4 mv_matrix;
uniform mat4 p_matrix;

attribute vec3 position;
attribute vec2 texcoord;

varying vec2 frag_texcoord;
varying float fade_factor;

void main()
{
    vec4 eye_position = mv_matrix * vec4(position, 1.0);
    gl_Position = p_matrix * eye_position;
    fade_factor = sin(timer) * 0.5 + 0.5;
    frag_texcoord = texcoord;
}
