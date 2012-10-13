#version 110

uniform sampler2D textures[2];

varying vec2 texcoord;
varying float fade_factor;

void main()
{
    gl_FragColor = mix(
        texture2D(textures[0], texcoord),
        texture2D(textures[1], texcoord),
        fade_factor
    );
}
