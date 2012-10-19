#version 120

uniform sampler2D textures[2];

varying vec2 frag_texcoord;
varying float fade_factor;

void main()
{
    gl_FragColor = mix(
        texture2D(textures[0], frag_texcoord),
        texture2D(textures[1], frag_texcoord),
        fade_factor
    );
}
