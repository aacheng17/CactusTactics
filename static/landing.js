import * as networking from './networking.js';
import { AVATARS } from './avatars/avatars.js';
import { initAudio } from './audio.js';
import { globalEn } from './globalEnum.js';

export var COLORS = ["chocolate", "crimson", "coral", "gold", "darkgreen", "springgreen", "turquoise", "cornflowerblue", "indigo", "orchid", "slategrey", "black"];

export var name;
var avatarIndex = 0;
var colorIndex = 0;
var landing = document.getElementById("landing");
var nameForm = document.getElementById("name-form");
var nameField = document.getElementById("name-field");
var avatarRandomize = document.getElementById("avatar-randomize");
var avatarButtonLeft = document.getElementById("avatar-button-left");
var avatarButtonRight = document.getElementById("avatar-button-right");
var avatarButtonColorLeft = document.getElementById("avatar-button-color-left");
var avatarButtonColorRight = document.getElementById("avatar-button-color-right");
var avatarSvg = document.getElementById("avatar-svg");
var avatarPath = document.getElementById("avatar-path");
var ingame = document.getElementById("ingame");

export function initLanding(conn) {
    nameForm.onsubmit = function (e) {
        initAudio();
        if (!conn) {
            return false;
        }
        if (!nameField.value.trim()) {
            return false;
        }
        name = nameField.value;
        networking.send(conn, globalEn.ToServerCode.NAME + name + "\t" + avatarIndex.toString() + "\t" + colorIndex);
        e.preventDefault();
        landing.parentNode.removeChild(landing);
        document.body.appendChild(ingame);
    }
    
    avatarRandomize.onclick = function(e) {
        randomizeAvatar();
    }

    avatarButtonLeft.onclick = function(e) {
        avatarIndex--;
        if (avatarIndex < 0) avatarIndex = AVATARS.length - 1;
        setAvatarSvg();
    }

    avatarButtonRight.onclick = function(e) {
        avatarIndex++;
        if (avatarIndex >= AVATARS.length) avatarIndex = 0;
        setAvatarSvg();
    }

    avatarButtonColorLeft.onclick = function(e) {
        colorIndex--;
        if (colorIndex < 0) colorIndex = COLORS.length - 1;
        avatarSvg.style.fill = COLORS[colorIndex];
    }

    avatarButtonColorRight.onclick = function(e) {
        colorIndex++;
        if (colorIndex >= COLORS.length) colorIndex = 0;
        avatarSvg.style.fill = COLORS[colorIndex];
    }

    randomizeAvatar();
}

function getRandomInt(max) {
    return Math.floor(Math.random() * max);
}

function randomizeAvatar() {
    avatarIndex = getRandomInt(AVATARS.length);
    setAvatarSvg();
    colorIndex = getRandomInt(COLORS.length);
    avatarSvg.style.fill = COLORS[colorIndex];
}

function setAvatarSvg() {
    avatarPath.setAttribute("d", AVATARS[avatarIndex]);
}